package balancecache

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
)

func (impl *Implementation) Invalidate(
	ctx context.Context,
	req dto.BalanceCacheInvalidateRequest,
) error {
	ctx, span := dspan.Start(
		ctx,
		"balancecache.Implementation.Invalidate",
		dspan.Attr("req", req),
	)
	defer span.End()

	if err := impl.client.Del(
		ctx,
		balanceKey(req.UserID.Value),
	).Err(); err != nil {
		return span.Err(err)
	}
	return nil
}
