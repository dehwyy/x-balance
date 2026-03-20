package transactionservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/application/service/transactionservice"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	transactionv1 "github.com/dehwyy/x-balance/internal/generated/pb/common/transaction/v1"
	"github.com/dehwyy/x-balance/pkg/storage"
	"github.com/dehwyy/x-balance/pkg/test/mocks"
)

const testUserID = "user-123"
const testEventID = "event-456"

func newTransactionService(eventRepo *mocks.EventRepository) *transactionservice.Service {
	return transactionservice.New(transactionservice.Opts{
		EventRepo: eventRepo,
	})
}

func TestListTransactions_Success(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}

	fromTime := time.Now().Add(-24 * time.Hour)
	toTime := time.Now()

	eventValues := []event.Event{
		{
			ID:            event.ID("event-1"),
			UserID:        user.ID(testUserID),
			Type:          transactionv1.TransactionType_TRANSACTION_TYPE_CREDIT,
			Amount:        event.Amount(decimal.NewFromInt(100)),
			TransactionID: event.TransactionID("tx-1"),
		},
		{
			ID:            event.ID("event-2"),
			UserID:        user.ID(testUserID),
			Type:          transactionv1.TransactionType_TRANSACTION_TYPE_DEBIT,
			Amount:        event.Amount(decimal.NewFromInt(-50)),
			TransactionID: event.TransactionID("tx-2"),
		},
	}

	expectedTotal := int64(2)

	expectedReq := dto.EventListRequest{
		UserID:     user.ID(testUserID),
		Pagination: storage.NewPagination(10, 0),
		From:       &fromTime,
		To:         &toTime,
	}

	eventRepo.On("List", mock.Anything, expectedReq).
		Return(dto.EventListResponse{Events: eventValues, Total: expectedTotal}, nil)

	svc := newTransactionService(eventRepo)
	resp, err := svc.ListTransactions(ctx, &transactionservice.ListTransactionsRequest{
		UserID:     user.ID(testUserID),
		Pagination: storage.NewPagination(10, 0),
		From:       &fromTime,
		To:         &toTime,
	})

	require.NoError(t, err)
	assert.Len(t, resp.Events, 2)
	assert.Equal(t, expectedTotal, resp.Total)
	eventRepo.AssertExpectations(t)
}

func TestListTransactions_EmptyResult(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}

	expectedReq := dto.EventListRequest{
		UserID:     user.ID(testUserID),
		Pagination: storage.NewPagination(10, 0),
	}

	eventRepo.On("List", mock.Anything, expectedReq).
		Return(dto.EventListResponse{Events: []event.Event{}, Total: int64(0)}, nil)

	svc := newTransactionService(eventRepo)
	resp, err := svc.ListTransactions(ctx, &transactionservice.ListTransactionsRequest{
		UserID:     user.ID(testUserID),
		Pagination: storage.NewPagination(10, 0),
	})

	require.NoError(t, err)
	assert.Empty(t, resp.Events)
	assert.Equal(t, int64(0), resp.Total)
	eventRepo.AssertExpectations(t)
}

func TestListTransactions_RepositoryError(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}

	expectedReq := dto.EventListRequest{
		UserID:     user.ID(testUserID),
		Pagination: storage.NewPagination(10, 0),
	}

	eventRepo.On("List", mock.Anything, expectedReq).
		Return(dto.EventListResponse{}, gorm.ErrInvalidDB)

	svc := newTransactionService(eventRepo)
	_, err := svc.ListTransactions(ctx, &transactionservice.ListTransactionsRequest{
		UserID:     user.ID(testUserID),
		Pagination: storage.NewPagination(10, 0),
	})

	assert.ErrorIs(t, err, gorm.ErrInvalidDB)
	eventRepo.AssertExpectations(t)
}

func TestGetTransaction_Success(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}

	expectedEvent := event.Event{
		ID:            event.ID(testEventID),
		UserID:        user.ID(testUserID),
		Type:          transactionv1.TransactionType_TRANSACTION_TYPE_CREDIT,
		Amount:        event.Amount(decimal.NewFromInt(100)),
		TransactionID: event.TransactionID("tx-1"),
	}

	eventRepo.On("GetByID", mock.Anything, dto.EventGetByIDRequest{ID: event.ID(testEventID)}).
		Return(dto.EventGetByIDResponse{Event: expectedEvent}, nil)

	svc := newTransactionService(eventRepo)
	resp, err := svc.GetTransaction(ctx, &transactionservice.GetTransactionRequest{
		UserID: user.ID(testUserID),
		TxID:   event.ID(testEventID),
	})

	require.NoError(t, err)
	assert.Equal(t, &expectedEvent, resp.Event)
	eventRepo.AssertExpectations(t)
}

func TestGetTransaction_NotFound(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}

	eventRepo.On("GetByID", mock.Anything, dto.EventGetByIDRequest{ID: event.ID(testEventID)}).
		Return(dto.EventGetByIDResponse{}, gorm.ErrRecordNotFound)

	svc := newTransactionService(eventRepo)
	_, err := svc.GetTransaction(ctx, &transactionservice.GetTransactionRequest{
		UserID: user.ID(testUserID),
		TxID:   event.ID(testEventID),
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	eventRepo.AssertExpectations(t)
}
