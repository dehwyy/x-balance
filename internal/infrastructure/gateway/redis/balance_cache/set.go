package balancecache

import (
	"context"
	"encoding/json"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
)

func (impl *Implementation) Set(
	ctx context.Context,
	req dto.BalanceCacheSetRequest,
) error {
	ctx, span := dspan.Start(
		ctx,
		"balancecache.Implementation.Set",
		dspan.Attr("req", req),
	)
	defer span.End()

	entry := cacheEntry{
		Available: req.Available.String(),
		Frozen:    req.Frozen.String(),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return span.Err(err)
	}

	if err := impl.client.Set(
		ctx,
		balanceKey(string(req.UserID)),
		data,
		0,
	).Err(); err != nil {
		return span.Err(err)
	}
	return nil
}
