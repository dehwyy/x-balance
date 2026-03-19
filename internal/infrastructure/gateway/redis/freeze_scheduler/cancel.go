package freezescheduler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
)

func (impl *Implementation) Cancel(
	ctx context.Context,
	req dto.FreezeCancelRequest,
) error {
	ctx, span := dspan.Start(ctx, "freezescheduler.Implementation.Cancel", dspan.Attr("req", req))
	defer span.End()

	if err := impl.client.Del(ctx, freezeKey(req.TransactionID.Value)).Err(); err != nil {
		return span.Err(err)
	}
	return nil
}
