package balanceservice_test

import (
	"context"
	"testing"

	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/pkg/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const testUserID = "user-1"

func makeSnap(balance string) *snapshot.Snapshot {
	b, _ := decimal.NewFromString(balance)
	return &snapshot.Snapshot{
		ID:      snapshot.ID{Value: "snap-1"},
		UserID:  testUserID,
		Balance: snapshot.Balance{Value: b},
		Version: snapshot.Version{Value: 1},
	}
}

func makeUser(overdraft string) *user.User {
	o, _ := decimal.NewFromString(overdraft)
	return &user.User{
		ID:             user.ID{Value: testUserID},
		OverdraftLimit: user.OverdraftLimit{Value: o},
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

	existingEvent := &event.Event{
		ID:            event.ID{Value: "ev-1"},
		UserID:        testUserID,
		TransactionID: event.TransactionID{Value: "tx-idempotent"},
		Amount:        event.Amount{Value: amount},
	}

	eventRepo.On("GetByTransactionID", ctx, event.TransactionID{Value: "tx-idempotent"}).Return(existingEvent, nil)
	snapshotRepo.On("GetLatestByUserID", ctx, testUserID).Return(snap, nil)
	zero := decimal.Zero
	eventRepo.On("SumSinceSnapshot", ctx, testUserID, snap.ID).Return(amount, zero, nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	resp, err := svc.Credit(ctx, balanceservice.CreditRequest{
		UserID:        testUserID,
		Amount:        amount,
		TransactionID: "tx-idempotent",
	})

	require.NoError(t, err)
	assert.Equal(t, "tx-idempotent", resp.TransactionID)
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

	eventRepo.On("GetByTransactionID", ctx, event.TransactionID{Value: "tx-1"}).Return(nil, gorm.ErrRecordNotFound)
	snapshotRepo.On("GetLatestByUserID", ctx, testUserID).Return(snap, nil)
	snapshotRepo.On("UpdateVersion", ctx, snap).Return(nil)

	createdEvent := &event.Event{
		ID:            event.ID{Value: "ev-new"},
		UserID:        testUserID,
		Type:          event.TypeCredit,
		Amount:        event.Amount{Value: amount},
		TransactionID: event.TransactionID{Value: "tx-1"},
	}
	eventRepo.On("Create", ctx, mock.AnythingOfType("*event.Event")).Return(createdEvent, nil)
	eventRepo.On("SumSinceSnapshot", ctx, testUserID, snap.ID).Return(amount, zero, nil)
	cache.On("Invalidate", ctx, testUserID).Return(nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	resp, err := svc.Credit(ctx, balanceservice.CreditRequest{
		UserID:        testUserID,
		Amount:        amount,
		TransactionID: "tx-1",
	})

	require.NoError(t, err)
	// balance = snap(100) + delta(50) - frozen(0) = 150
	assert.Equal(t, "150", resp.NewBalance.String())
	assert.Equal(t, "tx-1", resp.TransactionID)
}

func TestDebit_InsufficientFunds(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freeze := &mocks.FreezeScheduler{}

	snap := makeSnap("100")
	u := makeUser("0") // no overdraft
	debitAmount, _ := decimal.NewFromString("150")
	zero := decimal.Zero

	eventRepo.On("GetByTransactionID", ctx, event.TransactionID{Value: "tx-debit"}).Return(nil, gorm.ErrRecordNotFound)
	snapshotRepo.On("GetLatestByUserID", ctx, testUserID).Return(snap, nil)
	userRepo.On("GetByID", ctx, user.ID{Value: testUserID}).Return(u, nil)
	eventRepo.On("SumSinceSnapshot", ctx, testUserID, snap.ID).Return(zero, zero, nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	_, err := svc.Debit(ctx, balanceservice.DebitRequest{
		UserID:        testUserID,
		Amount:        debitAmount,
		TransactionID: "tx-debit",
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
	u := makeUser("50")                            // overdraft 50, so can go to -50
	debitAmount, _ := decimal.NewFromString("120") // 100 - 120 = -20, within -50 limit
	zero := decimal.Zero

	eventRepo.On("GetByTransactionID", ctx, event.TransactionID{Value: "tx-overdraft"}).Return(nil, gorm.ErrRecordNotFound)
	snapshotRepo.On("GetLatestByUserID", ctx, testUserID).Return(snap, nil)
	userRepo.On("GetByID", ctx, user.ID{Value: testUserID}).Return(u, nil)
	eventRepo.On("SumSinceSnapshot", ctx, testUserID, snap.ID).Return(zero, zero, nil)
	snapshotRepo.On("UpdateVersion", ctx, snap).Return(nil)
	eventRepo.On("Create", ctx, mock.AnythingOfType("*event.Event")).Return(&event.Event{
		UserID:        testUserID,
		Amount:        event.Amount{Value: debitAmount.Neg()},
		TransactionID: event.TransactionID{Value: "tx-overdraft"},
	}, nil)
	cache.On("Invalidate", ctx, testUserID).Return(nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	resp, err := svc.Debit(ctx, balanceservice.DebitRequest{
		UserID:        testUserID,
		Amount:        debitAmount,
		TransactionID: "tx-overdraft",
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
	debitAmount, _ := decimal.NewFromString("200") // 100 - 200 = -100, exceeds -50 limit
	zero := decimal.Zero

	eventRepo.On("GetByTransactionID", ctx, event.TransactionID{Value: "tx-exceed"}).Return(nil, gorm.ErrRecordNotFound)
	snapshotRepo.On("GetLatestByUserID", ctx, testUserID).Return(snap, nil)
	userRepo.On("GetByID", ctx, user.ID{Value: testUserID}).Return(u, nil)
	eventRepo.On("SumSinceSnapshot", ctx, testUserID, snap.ID).Return(zero, zero, nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	_, err := svc.Debit(ctx, balanceservice.DebitRequest{
		UserID:        testUserID,
		Amount:        debitAmount,
		TransactionID: "tx-exceed",
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

	eventRepo.On("GetByTransactionID", ctx, event.TransactionID{Value: "tx-freeze"}).Return(nil, gorm.ErrRecordNotFound)
	snapshotRepo.On("GetLatestByUserID", ctx, testUserID).Return(snap, nil)
	userRepo.On("GetByID", ctx, user.ID{Value: testUserID}).Return(u, nil)
	eventRepo.On("SumSinceSnapshot", ctx, testUserID, snap.ID).Return(zero, zero, nil)
	snapshotRepo.On("UpdateVersion", ctx, snap).Return(nil)
	eventRepo.On("Create", ctx, mock.AnythingOfType("*event.Event")).Return(&event.Event{
		UserID:        testUserID,
		Amount:        event.Amount{Value: freezeAmount},
		TransactionID: event.TransactionID{Value: "tx-freeze"},
	}, nil)
	cache.On("Invalidate", ctx, testUserID).Return(nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freezeSched)
	resp, err := svc.Freeze(ctx, balanceservice.FreezeRequest{
		UserID:               testUserID,
		Amount:               freezeAmount,
		TransactionID:        "tx-freeze",
		FreezeTimeoutSeconds: 0,
	})

	require.NoError(t, err)
	assert.Equal(t, "30", resp.FrozenAmount.String())
	assert.Equal(t, "tx-freeze", resp.TransactionID)
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

	eventRepo.On("GetByTransactionID", ctx, event.TransactionID{Value: "tx-freeze-ttl"}).Return(nil, gorm.ErrRecordNotFound)
	snapshotRepo.On("GetLatestByUserID", ctx, testUserID).Return(snap, nil)
	userRepo.On("GetByID", ctx, user.ID{Value: testUserID}).Return(u, nil)
	eventRepo.On("SumSinceSnapshot", ctx, testUserID, snap.ID).Return(zero, zero, nil)
	snapshotRepo.On("UpdateVersion", ctx, snap).Return(nil)
	eventRepo.On("Create", ctx, mock.AnythingOfType("*event.Event")).Return(&event.Event{
		UserID:        testUserID,
		Amount:        event.Amount{Value: freezeAmount},
		TransactionID: event.TransactionID{Value: "tx-freeze-ttl"},
	}, nil)
	freezeSched.On("Schedule", ctx, "tx-freeze-ttl", int64(10)).Return(nil)
	cache.On("Invalidate", ctx, testUserID).Return(nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freezeSched)
	resp, err := svc.Freeze(ctx, balanceservice.FreezeRequest{
		UserID:               testUserID,
		Amount:               freezeAmount,
		TransactionID:        "tx-freeze-ttl",
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

	freezeEvent := &event.Event{
		ID:            event.ID{Value: "ev-freeze"},
		UserID:        testUserID,
		Type:          event.TypeFreezeHold,
		Amount:        event.Amount{Value: freezeAmount},
		TransactionID: event.TransactionID{Value: "tx-freeze"},
	}

	releaseKey := "tx-freeze:release"

	eventRepo.On("GetByTransactionID", ctx, event.TransactionID{Value: releaseKey}).Return(nil, gorm.ErrRecordNotFound)
	eventRepo.On("GetByTransactionID", ctx, event.TransactionID{Value: "tx-freeze"}).Return(freezeEvent, nil)
	eventRepo.On("Create", ctx, mock.AnythingOfType("*event.Event")).Return(&event.Event{
		UserID:        testUserID,
		Amount:        event.Amount{Value: freezeAmount.Neg()},
		TransactionID: event.TransactionID{Value: releaseKey},
	}, nil)
	freezeSched.On("Cancel", ctx, "tx-freeze").Return(nil)
	cache.On("Invalidate", ctx, testUserID).Return(nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freezeSched)
	resp, err := svc.Unfreeze(ctx, balanceservice.UnfreezeRequest{
		UserID:        testUserID,
		TransactionID: "tx-freeze",
	})

	require.NoError(t, err)
	assert.Equal(t, "30", resp.UnfrozenAmount.String())
	assert.Equal(t, "tx-freeze", resp.TransactionID)
}

func TestUnfreeze_Idempotency(t *testing.T) {
	ctx := context.Background()

	eventRepo := &mocks.EventRepository{}
	snapshotRepo := &mocks.SnapshotRepository{}
	userRepo := &mocks.UserRepository{}
	cache := &mocks.BalanceCache{}
	freezeSched := &mocks.FreezeScheduler{}

	freezeAmount, _ := decimal.NewFromString("30")
	releaseKey := "tx-freeze:release"

	existingRelease := &event.Event{
		ID:            event.ID{Value: "ev-release"},
		UserID:        testUserID,
		Type:          event.TypeFreezeRelease,
		Amount:        event.Amount{Value: freezeAmount.Neg()},
		TransactionID: event.TransactionID{Value: releaseKey},
	}

	eventRepo.On("GetByTransactionID", ctx, event.TransactionID{Value: releaseKey}).Return(existingRelease, nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freezeSched)
	resp, err := svc.Unfreeze(ctx, balanceservice.UnfreezeRequest{
		UserID:        testUserID,
		TransactionID: "tx-freeze",
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

	cache.On("Get", ctx, testUserID).Return(available, frozen, true, nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	resp, err := svc.GetBalance(ctx, balanceservice.GetBalanceRequest{UserID: testUserID})

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

	cache.On("Get", ctx, testUserID).Return(zero, zero, false, nil)
	snapshotRepo.On("GetLatestByUserID", ctx, testUserID).Return(snap, nil)
	eventRepo.On("SumSinceSnapshot", ctx, testUserID, snap.ID).Return(delta, frozen, nil)
	cache.On("Set", ctx, testUserID, decimal.NewFromInt(130), frozen).Return(nil)

	svc := newService(eventRepo, snapshotRepo, userRepo, cache, freeze)
	resp, err := svc.GetBalance(ctx, balanceservice.GetBalanceRequest{UserID: testUserID})

	require.NoError(t, err)
	// available = 100 + 50 - 20 = 130
	assert.Equal(t, "130", resp.Available.String())
	assert.Equal(t, "20", resp.Frozen.String())
}
