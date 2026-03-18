package balanceservice

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/event"
)

func (s *Service) GetUserIDByTransactionID(
	ctx context.Context,
	txID string,
) (string, error) {
	e, err := s.eventRepo.GetByTransactionID(ctx, event.TransactionID{Value: txID})
	if err != nil {
		return "", err
	}
	return e.UserID, nil
}

func (s *Service) GetEventByTransactionID(
	ctx context.Context,
	txID string,
) (*event.Event, error) {
	return s.eventRepo.GetByTransactionID(ctx, event.TransactionID{Value: txID})
}
