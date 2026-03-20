package freezeexpiry

import (
	"context"
	"strings"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	tlog "github.com/dehwyy/tracerfx/pkg/tracer/log"

	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
)

func (w *Worker) handleExpiry(ctx context.Context, key string) {
	const prefix = "freeze:"
	if !strings.HasPrefix(key, prefix) {
		return
	}

	txID := strings.TrimPrefix(key, prefix)

	ctx, span := dspan.Start(
		ctx,
		"freezeexpiry.Worker.handleExpiry",
		dspan.Attr("tx_id", txID),
	)
	defer span.End()

	userIDResult, err := w.balanceservice.GetUserIDByTransactionID(
		ctx,
		&balanceservice.GetUserIDByTransactionIDRequest{
			TransactionID: event.NewTransactionID(txID),
		},
	)
	if err != nil {
		tlog.FromContext(ctx).Error("failed to find user for freeze expiry", "tx_id", txID, "err", err)
		return
	}

	if _, err := w.balanceservice.Unfreeze(
		ctx,
		&balanceservice.UnfreezeRequest{
			UserID:        userIDResult.UserID,
			TransactionID: event.NewTransactionID(txID),
		},
	); err != nil {
		tlog.FromContext(ctx).Error("failed to auto-unfreeze", "tx_id", txID, "err", err)
	}
}
