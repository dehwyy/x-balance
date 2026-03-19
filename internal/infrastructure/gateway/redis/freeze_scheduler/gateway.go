package freezescheduler

import (
	"fmt"

	"github.com/dehwyy/x-balance/internal/domain/gateway"
	"github.com/redis/go-redis/v9"
)

var _ gateway.FreezeScheduler = &Implementation{}

type Implementation struct {
	client *redis.Client
}

func New(client *redis.Client) *Implementation {
	return &Implementation{client: client}
}

func freezeKey(txID string) string {
	return fmt.Sprintf("freeze:%s", txID)
}
