package transactionservice

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/event"
)

type GetTransactionRequest struct {
	UserID string
	TxID   string
}

type GetTransactionResponse struct {
	Event *event.Event
}

func (s *Service) GetTransaction(
	ctx context.Context,
	req GetTransactionRequest,
) (*GetTransactionResponse, error) {
	e, err := s.eventRepo.GetByID(ctx, event.ID{Value: req.TxID})
	if err != nil {
		return nil, err
	}

	return &GetTransactionResponse{Event: e}, nil
}
