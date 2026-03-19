package transactionservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
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
	req *GetTransactionRequest,
) (*GetTransactionResponse, error) {
	ctx, span := dspan.Start(ctx, "transactionservice.Service.GetTransaction", dspan.Attr("req", req))
	defer span.End()

	getResp, err := s.eventRepo.GetByID(ctx, dto.EventGetByIDRequest{ID: event.ID{Value: req.TxID}})
	if err != nil {
		return nil, span.Err(err)
	}

	e := getResp.Event
	response := &GetTransactionResponse{Event: &e}
	span.WithAttribute("response", response)
	return response, nil
}
