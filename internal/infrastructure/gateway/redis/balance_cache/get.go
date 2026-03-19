package balancecache

import (
	"context"
	"encoding/json"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

func (impl *Implementation) Get(
	ctx context.Context,
	req dto.BalanceCacheGetRequest,
) (dto.BalanceCacheGetResponse, error) {
	ctx, span := dspan.Start(ctx, "balancecache.Implementation.Get", dspan.Attr("req", req))
	defer span.End()

	val, err := impl.client.Get(ctx, balanceKey(req.UserID.Value)).Result()
	if err != nil {
		if err == redis.Nil {
			response := dto.BalanceCacheGetResponse{Found: false}
			span.WithAttribute("response", response)
			return response, nil
		}
		return dto.BalanceCacheGetResponse{}, span.Err(err)
	}

	var entry cacheEntry
	if err := json.Unmarshal([]byte(val), &entry); err != nil {
		return dto.BalanceCacheGetResponse{}, span.Err(err)
	}

	available, err := decimal.NewFromString(entry.Available)
	if err != nil {
		return dto.BalanceCacheGetResponse{}, span.Err(err)
	}

	frozen, err := decimal.NewFromString(entry.Frozen)
	if err != nil {
		return dto.BalanceCacheGetResponse{}, span.Err(err)
	}

	response := dto.BalanceCacheGetResponse{Available: available, Frozen: frozen, Found: true}
	span.WithAttribute("response", response)
	return response, nil
}
