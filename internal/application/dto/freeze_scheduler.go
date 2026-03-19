package dto

import (
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
)

type FreezeScheduleRequest struct {
	TransactionID event.TransactionID
	TTLSeconds    int64
}

type FreezeCancelRequest struct {
	TransactionID event.TransactionID
}
