package transactionservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type GetTransactionRequest struct {
	UserID user.ID
	TxID   event.ID
}

type GetTransactionResponse struct {
	Event *event.Event
}

func (s *Service) GetTransaction(
	ctx context.Context,
	req *GetTransactionRequest,
) (*GetTransactionResponse, error) {
	ctx, span := dspan.Start(
		ctx,
		"transactionservice.Service.GetTransaction",
		dspan.Attr("req", req),
	)
	defer span.End()

	eventResult, err := s.eventRepo.GetByID(
		ctx,
		dto.EventGetByIDRequest{ID: req.TxID},
	)
	if err != nil {
		return nil, span.Err(err)
	}

	e := eventResult.Event
	return dspan.Response(span, &GetTransactionResponse{Event: &e}), nil
}
