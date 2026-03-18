package worker

import (
	"context"
	"strings"

	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type FreezeExpiryWorker struct {
	client     *redis.Client
	balanceSvc *balanceservice.Service
	log        zerolog.Logger
}

func NewFreezeExpiryWorker(
	client *redis.Client,
	balanceSvc *balanceservice.Service,
	log zerolog.Logger,
) *FreezeExpiryWorker {
	return &FreezeExpiryWorker{
		client:     client,
		balanceSvc: balanceSvc,
		log:        log,
	}
}

func (w *FreezeExpiryWorker) Start(ctx context.Context) {
	pubsub := w.client.Subscribe(ctx, "__keyevent@0__:expired")
	defer func() { _ = pubsub.Close() }()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			w.handleExpiry(ctx, msg.Payload)
		}
	}
}

func (w *FreezeExpiryWorker) handleExpiry(ctx context.Context, key string) {
	const prefix = "freeze:"
	if !strings.HasPrefix(key, prefix) {
		return
	}

	txID := strings.TrimPrefix(key, prefix)

	userID, err := w.balanceSvc.GetUserIDByTransactionID(ctx, txID)
	if err != nil {
		w.log.Error().Err(err).Str("tx_id", txID).Msg("failed to find user for freeze expiry")
		return
	}

	_, err = w.balanceSvc.Unfreeze(ctx, balanceservice.UnfreezeRequest{
		UserID:        userID,
		TransactionID: txID,
	})
	if err != nil {
		w.log.Error().Err(err).Str("tx_id", txID).Msg("failed to auto-unfreeze")
	}
}
