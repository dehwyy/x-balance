package balanceservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type GetUserIDByTransactionIDRequest struct {
	TransactionID event.TransactionID
}

type GetUserIDByTransactionIDResponse struct {
	UserID user.ID
}

type GetEventByTransactionIDRequest struct {
	TransactionID event.TransactionID
}

type GetEventByTransactionIDResponse struct {
	Event event.Event
}

func (s *Service) GetUserIDByTransactionID(
	ctx context.Context,
	req *GetUserIDByTransactionIDRequest,
) (*GetUserIDByTransactionIDResponse, error) {
	ctx, span := dspan.Start(
		ctx,
		"balanceservice.Service.GetUserIDByTransactionID",
		dspan.Attr("req", req),
	)
	defer span.End()

	eventDTO, err := s.eventRepo.GetByTransactionID(
		ctx,
		dto.EventGetByTxIDRequest{TransactionID: req.TransactionID},
	)
	if err != nil {
		return nil, span.Err(err)
	}

	return dspan.Response(span, &GetUserIDByTransactionIDResponse{UserID: eventDTO.Event.UserID}), nil
}

func (s *Service) GetEventByTransactionID(
	ctx context.Context,
	req *GetEventByTransactionIDRequest,
) (*GetEventByTransactionIDResponse, error) {
	ctx, span := dspan.Start(
		ctx,
		"balanceservice.Service.GetEventByTransactionID",
		dspan.Attr("req", req),
	)
	defer span.End()

	eventDTO, err := s.eventRepo.GetByTransactionID(
		ctx,
		dto.EventGetByTxIDRequest{TransactionID: req.TransactionID},
	)
	if err != nil {
		return nil, span.Err(err)
	}

	return dspan.Response(span, &GetEventByTransactionIDResponse{Event: eventDTO.Event}), nil
}
