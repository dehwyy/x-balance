package transactionservice

import (
	"context"
	"time"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/pkg/storage"
)

type ListTransactionsRequest struct {
	UserID     string
	Pagination storage.Pagination
	From       *time.Time
	To         *time.Time
}

type ListTransactionsResponse struct {
	Events []*event.Event
	Total  int64
}

func (s *Service) ListTransactions(
	ctx context.Context,
	req *ListTransactionsRequest,
) (*ListTransactionsResponse, error) {
	ctx, span := dspan.Start(ctx, "transactionservice.Service.ListTransactions", dspan.Attr("req", req))
	defer span.End()

	listResp, err := s.eventRepo.List(ctx, dto.EventListRequest{
		UserID:     user.ID{Value: req.UserID},
		Pagination: req.Pagination,
		From:       req.From,
		To:         req.To,
	})
	if err != nil {
		return nil, span.Err(err)
	}

	events := make([]*event.Event, len(listResp.Events))
	for i := range listResp.Events {
		e := listResp.Events[i]
		events[i] = &e
	}

	response := &ListTransactionsResponse{Events: events, Total: listResp.Total}
	span.WithAttribute("response", response)
	return response, nil
}
