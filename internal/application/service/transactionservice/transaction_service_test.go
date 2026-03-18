package transactionservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/dehwyy/x-balance/internal/application/service/transactionservice"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/dehwyy/x-balance/pkg/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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

	expectedEvents := []*event.Event{
		{
			ID:            event.ID{Value: "event-1"},
			UserID:        testUserID,
			Type:          event.TypeCredit,
			Amount:        event.Amount{Value: decimal.NewFromInt(100)},
			TransactionID: event.TransactionID{Value: "tx-1"},
		},
		{
			ID:            event.ID{Value: "event-2"},
			UserID:        testUserID,
			Type:          event.TypeDebit,
			Amount:        event.Amount{Value: decimal.NewFromInt(-50)},
			TransactionID: event.TransactionID{Value: "tx-2"},
		},
	}

	expectedTotal := int64(2)

	expectedReq := repository.ListEventsRequest{
		UserID: testUserID,
		Limit:  10,
		Offset: 0,
		From:   &fromTime,
		To:     &toTime,
	}

	eventRepo.On("List", ctx, expectedReq).Return(expectedEvents, expectedTotal, nil)

	svc := newTransactionService(eventRepo)
	resp, err := svc.ListTransactions(ctx, transactionservice.ListTransactionsRequest{
		UserID: testUserID,
		Limit:  10,
		Offset: 0,
		From:   &fromTime,
		To:     &toTime,
	})

	require.NoError(t, err)
	assert.Equal(t, expectedEvents, resp.Events)
	assert.Equal(t, expectedTotal, resp.Total)
	eventRepo.AssertExpectations(t)
}

func TestListTransactions_EmptyResult(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}

	expectedEvents := []*event.Event{}
	expectedTotal := int64(0)

	expectedReq := repository.ListEventsRequest{
		UserID: testUserID,
		Limit:  10,
		Offset: 0,
		From:   nil,
		To:     nil,
	}

	eventRepo.On("List", ctx, expectedReq).Return(expectedEvents, expectedTotal, nil)

	svc := newTransactionService(eventRepo)
	resp, err := svc.ListTransactions(ctx, transactionservice.ListTransactionsRequest{
		UserID: testUserID,
		Limit:  10,
		Offset: 0,
	})

	require.NoError(t, err)
	assert.Empty(t, resp.Events)
	assert.Equal(t, int64(0), resp.Total)
	eventRepo.AssertExpectations(t)
}

func TestListTransactions_RepositoryError(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}

	expectedReq := repository.ListEventsRequest{
		UserID: testUserID,
		Limit:  10,
		Offset: 0,
		From:   nil,
		To:     nil,
	}

	eventRepo.On("List", ctx, expectedReq).Return(nil, int64(0), gorm.ErrInvalidDB)

	svc := newTransactionService(eventRepo)
	_, err := svc.ListTransactions(ctx, transactionservice.ListTransactionsRequest{
		UserID: testUserID,
		Limit:  10,
		Offset: 0,
	})

	assert.ErrorIs(t, err, gorm.ErrInvalidDB)
	eventRepo.AssertExpectations(t)
}

func TestGetTransaction_Success(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}

	expectedEvent := &event.Event{
		ID:            event.ID{Value: testEventID},
		UserID:        testUserID,
		Type:          event.TypeCredit,
		Amount:        event.Amount{Value: decimal.NewFromInt(100)},
		TransactionID: event.TransactionID{Value: "tx-1"},
	}

	eventRepo.On("GetByID", ctx, event.ID{Value: testEventID}).Return(expectedEvent, nil)

	svc := newTransactionService(eventRepo)
	resp, err := svc.GetTransaction(ctx, transactionservice.GetTransactionRequest{
		UserID: testUserID,
		TxID:   testEventID,
	})

	require.NoError(t, err)
	assert.Equal(t, expectedEvent, resp.Event)
	eventRepo.AssertExpectations(t)
}

func TestGetTransaction_NotFound(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}

	eventRepo.On("GetByID", ctx, event.ID{Value: testEventID}).Return(nil, gorm.ErrRecordNotFound)

	svc := newTransactionService(eventRepo)
	_, err := svc.GetTransaction(ctx, transactionservice.GetTransactionRequest{
		UserID: testUserID,
		TxID:   testEventID,
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	eventRepo.AssertExpectations(t)
}
