package transactionservice

import (
	"context"
	"time"

	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/domain/repository"
)

type ListTransactionsRequest struct {
	UserID string
	Limit  int
	Offset int
	From   *time.Time
	To     *time.Time
}

type ListTransactionsResponse struct {
	Events []*event.Event
	Total  int64
}

func (s *Service) ListTransactions(
	ctx context.Context,
	req ListTransactionsRequest,
) (*ListTransactionsResponse, error) {
	events, total, err := s.eventRepo.List(ctx, repository.ListEventsRequest{
		UserID: req.UserID,
		Limit:  req.Limit,
		Offset: req.Offset,
		From:   req.From,
		To:     req.To,
	})
	if err != nil {
		return nil, err
	}

	return &ListTransactionsResponse{Events: events, Total: total}, nil
}
