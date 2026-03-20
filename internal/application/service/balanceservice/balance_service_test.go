package balanceservice_test

import (
	"context"
	"testing"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	transactionv1 "github.com/dehwyy/x-balance/internal/generated/pb/common/transaction/v1"
	"github.com/dehwyy/x-balance/pkg/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const testUserID = "user-1"

func makeSnap(balance string) *snapshot.Snapshot {
	b, _ := decimal.NewFromString(balance)
	return &snapshot.Snapshot{
		ID:      snapshot.ID("snap-1"),
		UserID:  user.ID(testUserID),
		Balance: snapshot.Balance(b),
		Version: snapshot.Version(1),
	}
}

func makeUser(overdraft string) *user.User {
	o, _ := decimal.NewFromString(overdraft)
	return &user.User{
		ID:             user.ID(testUserID),
		OverdraftLimit: user.OverdraftLimit(o),
	}
}

func newService(
	eventRepo *mocks.EventRepository,
	snapshotRepo *mocks.SnapshotRepository,
	userRepo *mocks.UserRepository,
	cache *mocks.BalanceCache,
	freezeScheduler *mocks.FreezeScheduler,
) *balanceservice.Service {
	tx := &mocks.TxManager{}
	return balanceservice.New(balanceservice.Opts{
		TX:              tx,
		EventRepo:       eventRepo,
		SnapshotRepo:    snapshotRepo,
		UserRepo:        userRepo,
		BalanceCache:    cache,
		FreezeScheduler: freezeScheduler,
		Config:          balanceservice.BalanceConfig{SnapshotEveryN: 0},
	})
}

func TestCredit_Idempotency(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freeze := &mocks.FreezeScheduler{}

	snap := makeSnap("100")
	amount, _ := decimal.NewFromString("50")
	zero := decimal.Zero

	existingEvent := event.Event{
		ID:            event.ID("ev-1"),
		UserID:        user.ID(testUserID),
		TransactionID: event.TransactionID("tx-idempotent"),
		Amount:        event.Amount(amount),
	}

	eventRepo.On("GetByTransactionID", mock.Anything, dto.EventGetByTxIDRequest{TransactionID: event.TransactionID("tx-idempotent")}).
		Return(dto.EventGetByTxIDResponse{Event: existingEvent}, nil)
	snapshotRepo.On("GetLatestByUserID", mock.Anything, dto.SnapshotGetLatestByUserIDRequest{UserID: user.ID(testUserID)}).
		Return(dto.SnapshotGetLatestByUserIDResponse{Snapshot: *snap}, nil)
	eventRepo.On("SumSinceSnapshot", mock.Anything, dto.EventSumSinceSnapshotRequest{UserID: user.ID(testUserID), SnapshotID: snap.ID}).
		Return(dto.EventSumSinceSnapshotResponse{Available: amount, Frozen: zero}, nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	resp, err := svc.Credit(ctx, &balanceservice.CreditRequest{
		UserID:        user.ID(testUserID),
		Amount:        amount,
		TransactionID: event.TransactionID("tx-idempotent"),
	})

	require.NoError(t, err)
	assert.Equal(t, "tx-idempotent", string(resp.TransactionID))
	eventRepo.AssertExpectations(t)
}

func TestCredit_Success(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freeze := &mocks.FreezeScheduler{}

	snap := makeSnap("100")
	amount, _ := decimal.NewFromString("50")
	zero := decimal.Zero

	eventRepo.On("GetByTransactionID", mock.Anything, dto.EventGetByTxIDRequest{TransactionID: event.TransactionID("tx-1")}).
		Return(dto.EventGetByTxIDResponse{}, repository.ErrNotFound)
	snapshotRepo.On("GetLatestByUserID", mock.Anything, dto.SnapshotGetLatestByUserIDRequest{UserID: user.ID(testUserID)}).
		Return(dto.SnapshotGetLatestByUserIDResponse{Snapshot: *snap}, nil)
	snapshotRepo.On("UpdateVersion", mock.Anything, dto.SnapshotUpdateVersionRequest{Snapshot: *snap}).
		Return(nil)
	eventRepo.On("Create", mock.Anything, mock.AnythingOfType("dto.EventCreateRequest")).
		Return(dto.EventCreateResponse{Event: event.Event{
			ID:            event.ID("ev-new"),
			UserID:        user.ID(testUserID),
			Type:          transactionv1.TransactionType_TRANSACTION_TYPE_CREDIT,
			Amount:        event.Amount(amount),
			TransactionID: event.TransactionID("tx-1"),
		}}, nil)
	eventRepo.On("SumSinceSnapshot", mock.Anything, dto.EventSumSinceSnapshotRequest{UserID: user.ID(testUserID), SnapshotID: snap.ID}).
		Return(dto.EventSumSinceSnapshotResponse{Available: amount, Frozen: zero}, nil)
	cache.On("Invalidate", mock.Anything, dto.BalanceCacheInvalidateRequest{UserID: user.ID(testUserID)}).
		Return(nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	resp, err := svc.Credit(ctx, &balanceservice.CreditRequest{
		UserID:        user.ID(testUserID),
		Amount:        amount,
		TransactionID: event.TransactionID("tx-1"),
	})

	require.NoError(t, err)
	assert.Equal(t, "150", resp.NewBalance.String())
	assert.Equal(t, "tx-1", string(resp.TransactionID))
}

func TestDebit_InsufficientFunds(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freeze := &mocks.FreezeScheduler{}

	snap := makeSnap("100")
	u := makeUser("0")
	debitAmount, _ := decimal.NewFromString("150")
	zero := decimal.Zero

	eventRepo.On("GetByTransactionID", mock.Anything, dto.EventGetByTxIDRequest{TransactionID: event.TransactionID("tx-debit")}).
		Return(dto.EventGetByTxIDResponse{}, repository.ErrNotFound)
	snapshotRepo.On("GetLatestByUserID", mock.Anything, dto.SnapshotGetLatestByUserIDRequest{UserID: user.ID(testUserID)}).
		Return(dto.SnapshotGetLatestByUserIDResponse{Snapshot: *snap}, nil)
	userRepo.On("GetByID", mock.Anything, dto.UserGetByIDRequest{ID: user.ID(testUserID)}).
		Return(dto.UserGetByIDResponse{User: *u}, nil)
	eventRepo.On("SumSinceSnapshot", mock.Anything, dto.EventSumSinceSnapshotRequest{UserID: user.ID(testUserID), SnapshotID: snap.ID}).
		Return(dto.EventSumSinceSnapshotResponse{Available: zero, Frozen: zero}, nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	_, err := svc.Debit(ctx, &balanceservice.DebitRequest{
		UserID:        user.ID(testUserID),
		Amount:        debitAmount,
		TransactionID: event.TransactionID("tx-debit"),
	})

	assert.ErrorIs(t, err, balanceservice.ErrInsufficientFunds)
}

func TestDebit_WithinOverdraft(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freeze := &mocks.FreezeScheduler{}

	snap := makeSnap("100")
	u := makeUser("50")
	debitAmount, _ := decimal.NewFromString("120")
	zero := decimal.Zero

	eventRepo.On("GetByTransactionID", mock.Anything, dto.EventGetByTxIDRequest{TransactionID: event.TransactionID("tx-overdraft")}).
		Return(dto.EventGetByTxIDResponse{}, repository.ErrNotFound)
	snapshotRepo.On("GetLatestByUserID", mock.Anything, dto.SnapshotGetLatestByUserIDRequest{UserID: user.ID(testUserID)}).
		Return(dto.SnapshotGetLatestByUserIDResponse{Snapshot: *snap}, nil)
	userRepo.On("GetByID", mock.Anything, dto.UserGetByIDRequest{ID: user.ID(testUserID)}).
		Return(dto.UserGetByIDResponse{User: *u}, nil)
	eventRepo.On("SumSinceSnapshot", mock.Anything, dto.EventSumSinceSnapshotRequest{UserID: user.ID(testUserID), SnapshotID: snap.ID}).
		Return(dto.EventSumSinceSnapshotResponse{Available: zero, Frozen: zero}, nil)
	snapshotRepo.On("UpdateVersion", mock.Anything, dto.SnapshotUpdateVersionRequest{Snapshot: *snap}).
		Return(nil)
	eventRepo.On("Create", mock.Anything, mock.AnythingOfType("dto.EventCreateRequest")).
		Return(dto.EventCreateResponse{Event: event.Event{
			UserID:        user.ID(testUserID),
			Amount:        event.Amount(debitAmount.Neg()),
			TransactionID: event.TransactionID("tx-overdraft"),
		}}, nil)
	cache.On("Invalidate", mock.Anything, dto.BalanceCacheInvalidateRequest{UserID: user.ID(testUserID)}).
		Return(nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	resp, err := svc.Debit(ctx, &balanceservice.DebitRequest{
		UserID:        user.ID(testUserID),
		Amount:        debitAmount,
		TransactionID: event.TransactionID("tx-overdraft"),
	})

	require.NoError(t, err)
	assert.Equal(t, "-20", resp.NewBalance.String())
}

func TestDebit_ExceedsOverdraft(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freeze := &mocks.FreezeScheduler{}

	snap := makeSnap("100")
	u := makeUser("50")
	debitAmount, _ := decimal.NewFromString("200")
	zero := decimal.Zero

	eventRepo.On("GetByTransactionID", mock.Anything, dto.EventGetByTxIDRequest{TransactionID: event.TransactionID("tx-exceed")}).
		Return(dto.EventGetByTxIDResponse{}, repository.ErrNotFound)
	snapshotRepo.On("GetLatestByUserID", mock.Anything, dto.SnapshotGetLatestByUserIDRequest{UserID: user.ID(testUserID)}).
		Return(dto.SnapshotGetLatestByUserIDResponse{Snapshot: *snap}, nil)
	userRepo.On("GetByID", mock.Anything, dto.UserGetByIDRequest{ID: user.ID(testUserID)}).
		Return(dto.UserGetByIDResponse{User: *u}, nil)
	eventRepo.On("SumSinceSnapshot", mock.Anything, dto.EventSumSinceSnapshotRequest{UserID: user.ID(testUserID), SnapshotID: snap.ID}).
		Return(dto.EventSumSinceSnapshotResponse{Available: zero, Frozen: zero}, nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	_, err := svc.Debit(ctx, &balanceservice.DebitRequest{
		UserID:        user.ID(testUserID),
		Amount:        debitAmount,
		TransactionID: event.TransactionID("tx-exceed"),
	})

	assert.ErrorIs(t, err, balanceservice.ErrInsufficientFunds)
}

func TestFreeze_Success(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freezeSched := &mocks.FreezeScheduler{}

	snap := makeSnap("100")
	u := makeUser("0")
	freezeAmount, _ := decimal.NewFromString("30")
	zero := decimal.Zero

	eventRepo.On("GetByTransactionID", mock.Anything, dto.EventGetByTxIDRequest{TransactionID: event.TransactionID("tx-freeze")}).
		Return(dto.EventGetByTxIDResponse{}, repository.ErrNotFound)
	snapshotRepo.On("GetLatestByUserID", mock.Anything, dto.SnapshotGetLatestByUserIDRequest{UserID: user.ID(testUserID)}).
		Return(dto.SnapshotGetLatestByUserIDResponse{Snapshot: *snap}, nil)
	userRepo.On("GetByID", mock.Anything, dto.UserGetByIDRequest{ID: user.ID(testUserID)}).
		Return(dto.UserGetByIDResponse{User: *u}, nil)
	eventRepo.On("SumSinceSnapshot", mock.Anything, dto.EventSumSinceSnapshotRequest{UserID: user.ID(testUserID), SnapshotID: snap.ID}).
		Return(dto.EventSumSinceSnapshotResponse{Available: zero, Frozen: zero}, nil)
	snapshotRepo.On("UpdateVersion", mock.Anything, dto.SnapshotUpdateVersionRequest{Snapshot: *snap}).
		Return(nil)
	eventRepo.On("Create", mock.Anything, mock.AnythingOfType("dto.EventCreateRequest")).
		Return(dto.EventCreateResponse{Event: event.Event{
			UserID:        user.ID(testUserID),
			Amount:        event.Amount(freezeAmount),
			TransactionID: event.TransactionID("tx-freeze"),
		}}, nil)
	cache.On("Invalidate", mock.Anything, dto.BalanceCacheInvalidateRequest{UserID: user.ID(testUserID)}).
		Return(nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freezeSched)
	resp, err := svc.Freeze(ctx, &balanceservice.FreezeRequest{
		UserID:               user.ID(testUserID),
		Amount:               freezeAmount,
		TransactionID:        event.TransactionID("tx-freeze"),
		FreezeTimeoutSeconds: 0,
	})

	require.NoError(t, err)
	assert.Equal(t, "30", resp.FrozenAmount.String())
	assert.Equal(t, "tx-freeze", string(resp.TransactionID))
}

func TestFreeze_WithTimeout(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freezeSched := &mocks.FreezeScheduler{}

	snap := makeSnap("100")
	u := makeUser("0")
	freezeAmount, _ := decimal.NewFromString("30")
	zero := decimal.Zero

	eventRepo.On("GetByTransactionID", mock.Anything, dto.EventGetByTxIDRequest{TransactionID: event.TransactionID("tx-freeze-ttl")}).
		Return(dto.EventGetByTxIDResponse{}, repository.ErrNotFound)
	snapshotRepo.On("GetLatestByUserID", mock.Anything, dto.SnapshotGetLatestByUserIDRequest{UserID: user.ID(testUserID)}).
		Return(dto.SnapshotGetLatestByUserIDResponse{Snapshot: *snap}, nil)
	userRepo.On("GetByID", mock.Anything, dto.UserGetByIDRequest{ID: user.ID(testUserID)}).
		Return(dto.UserGetByIDResponse{User: *u}, nil)
	eventRepo.On("SumSinceSnapshot", mock.Anything, dto.EventSumSinceSnapshotRequest{UserID: user.ID(testUserID), SnapshotID: snap.ID}).
		Return(dto.EventSumSinceSnapshotResponse{Available: zero, Frozen: zero}, nil)
	snapshotRepo.On("UpdateVersion", mock.Anything, dto.SnapshotUpdateVersionRequest{Snapshot: *snap}).
		Return(nil)
	eventRepo.On("Create", mock.Anything, mock.AnythingOfType("dto.EventCreateRequest")).
		Return(dto.EventCreateResponse{Event: event.Event{
			UserID:        user.ID(testUserID),
			Amount:        event.Amount(freezeAmount),
			TransactionID: event.TransactionID("tx-freeze-ttl"),
		}}, nil)
	freezeSched.On("Schedule", mock.Anything, dto.FreezeScheduleRequest{TransactionID: event.TransactionID("tx-freeze-ttl"), TTLSeconds: int64(10)}).
		Return(nil)
	cache.On("Invalidate", mock.Anything, dto.BalanceCacheInvalidateRequest{UserID: user.ID(testUserID)}).
		Return(nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freezeSched)
	resp, err := svc.Freeze(ctx, &balanceservice.FreezeRequest{
		UserID:               user.ID(testUserID),
		Amount:               freezeAmount,
		TransactionID:        event.TransactionID("tx-freeze-ttl"),
		FreezeTimeoutSeconds: 10,
	})

	require.NoError(t, err)
	assert.Equal(t, "30", resp.FrozenAmount.String())
	freezeSched.AssertExpectations(t)
}

func TestUnfreeze_Success(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freezeSched := &mocks.FreezeScheduler{}

	freezeAmount, _ := decimal.NewFromString("30")

	freezeEvent := event.Event{
		ID:            event.ID("ev-freeze"),
		UserID:        user.ID(testUserID),
		Type:          transactionv1.TransactionType_TRANSACTION_TYPE_FREEZE_HOLD,
		Amount:        event.Amount(freezeAmount),
		TransactionID: event.TransactionID("tx-freeze"),
	}

	releaseKey := event.TransactionID("tx-freeze:release")

	eventRepo.On("GetByTransactionID", mock.Anything, dto.EventGetByTxIDRequest{TransactionID: releaseKey}).
		Return(dto.EventGetByTxIDResponse{}, repository.ErrNotFound)
	eventRepo.On("GetByTransactionID", mock.Anything, dto.EventGetByTxIDRequest{TransactionID: event.TransactionID("tx-freeze")}).
		Return(dto.EventGetByTxIDResponse{Event: freezeEvent}, nil)
	eventRepo.On("Create", mock.Anything, mock.AnythingOfType("dto.EventCreateRequest")).
		Return(dto.EventCreateResponse{Event: event.Event{
			UserID:        user.ID(testUserID),
			Amount:        event.Amount(freezeAmount.Neg()),
			TransactionID: releaseKey,
		}}, nil)
	freezeSched.On("Cancel", mock.Anything, dto.FreezeCancelRequest{TransactionID: event.TransactionID("tx-freeze")}).
		Return(nil)
	cache.On("Invalidate", mock.Anything, dto.BalanceCacheInvalidateRequest{UserID: user.ID(testUserID)}).
		Return(nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freezeSched)
	resp, err := svc.Unfreeze(ctx, &balanceservice.UnfreezeRequest{
		UserID:        user.ID(testUserID),
		TransactionID: event.TransactionID("tx-freeze"),
	})

	require.NoError(t, err)
	assert.Equal(t, "30", resp.UnfrozenAmount.String())
	assert.Equal(t, "tx-freeze", string(resp.TransactionID))
}

func TestUnfreeze_Idempotency(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freezeSched := &mocks.FreezeScheduler{}

	freezeAmount, _ := decimal.NewFromString("30")
	releaseKey := event.TransactionID("tx-freeze:release")

	existingRelease := event.Event{
		ID:            event.ID("ev-release"),
		UserID:        user.ID(testUserID),
		Type:          transactionv1.TransactionType_TRANSACTION_TYPE_FREEZE_RELEASE,
		Amount:        event.Amount(freezeAmount.Neg()),
		TransactionID: releaseKey,
	}

	eventRepo.On("GetByTransactionID", mock.Anything, dto.EventGetByTxIDRequest{TransactionID: releaseKey}).
		Return(dto.EventGetByTxIDResponse{Event: existingRelease}, nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freezeSched)
	resp, err := svc.Unfreeze(ctx, &balanceservice.UnfreezeRequest{
		UserID:        user.ID(testUserID),
		TransactionID: event.TransactionID("tx-freeze"),
	})

	require.NoError(t, err)
	assert.Equal(t, "30", resp.UnfrozenAmount.String())
	eventRepo.AssertNumberOfCalls(t, "Create", 0)
}

func TestGetBalance_CacheHit(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freeze := &mocks.FreezeScheduler{}

	available, _ := decimal.NewFromString("200")
	frozen, _ := decimal.NewFromString("50")

	cache.On("Get", mock.Anything, dto.BalanceCacheGetRequest{UserID: user.ID(testUserID)}).
		Return(dto.BalanceCacheGetResponse{Available: available, Frozen: frozen, Found: true}, nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	resp, err := svc.GetBalance(ctx, &balanceservice.GetBalanceRequest{UserID: user.ID(testUserID)})

	require.NoError(t, err)
	assert.Equal(t, "200", resp.Available.String())
	assert.Equal(t, "50", resp.Frozen.String())
	assert.Equal(t, "250", resp.Total.String())
}

func TestGetBalance_CacheMiss(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freeze := &mocks.FreezeScheduler{}

	snap := makeSnap("100")
	delta, _ := decimal.NewFromString("50")
	frozen, _ := decimal.NewFromString("20")
	zero := decimal.Zero

	cache.On("Get", mock.Anything, dto.BalanceCacheGetRequest{UserID: user.ID(testUserID)}).
		Return(dto.BalanceCacheGetResponse{Available: zero, Frozen: zero, Found: false}, nil)
	snapshotRepo.On("GetLatestByUserID", mock.Anything, dto.SnapshotGetLatestByUserIDRequest{UserID: user.ID(testUserID)}).
		Return(dto.SnapshotGetLatestByUserIDResponse{Snapshot: *snap}, nil)
	eventRepo.On("SumSinceSnapshot", mock.Anything, dto.EventSumSinceSnapshotRequest{UserID: user.ID(testUserID), SnapshotID: snap.ID}).
		Return(dto.EventSumSinceSnapshotResponse{Available: delta, Frozen: frozen}, nil)
	cache.On("Set", mock.Anything, dto.BalanceCacheSetRequest{UserID: user.ID(testUserID), Available: decimal.NewFromInt(130), Frozen: frozen}).
		Return(nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	resp, err := svc.GetBalance(ctx, &balanceservice.GetBalanceRequest{UserID: user.ID(testUserID)})

	require.NoError(t, err)
	assert.Equal(t, "130", resp.Available.String())
	assert.Equal(t, "20", resp.Frozen.String())
}
