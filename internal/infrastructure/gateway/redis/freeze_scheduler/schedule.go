package freezescheduler

import (
	"context"
	"time"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
)

func (impl *Implementation) Schedule(
	ctx context.Context,
	req dto.FreezeScheduleRequest,
) error {
	ctx, span := dspan.Start(ctx, "freezescheduler.Implementation.Schedule", dspan.Attr("req", req))
	defer span.End()

	if req.TTLSeconds <= 0 {
		return nil
	}

	if err := impl.client.Set(
		ctx,
		freezeKey(req.TransactionID.Value),
		"",
		time.Duration(req.TTLSeconds)*time.Second,
	).Err(); err != nil {
		return span.Err(err)
	}
	return nil
}
