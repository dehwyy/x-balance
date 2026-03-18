package balancehandler

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	"github.com/dehwyy/x-balance/internal/delivery/api/balance/convert"
	balancepb "github.com/dehwyy/x-balance/internal/generated/pb"
)

type Handler struct {
	balancepb.UnimplementedBalanceServiceServer
	svc *balanceservice.Service
}

func New(svc *balanceservice.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetBalance(
	ctx context.Context,
	req *balancepb.GetBalanceRequest,
) (*balancepb.GetBalanceResponse, error) {
	res, err := h.svc.GetBalance(ctx, balanceservice.GetBalanceRequest{UserID: req.UserId})
	if err != nil {
		return nil, err
	}
	return convert.GetBalanceResponseToProto(res), nil
}

func (h *Handler) Credit(
	ctx context.Context,
	req *balancepb.CreditRequest,
) (*balancepb.CreditResponse, error) {
	res, err := h.svc.Credit(ctx, convert.CreditRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return convert.CreditResponseToProto(res), nil
}

func (h *Handler) Debit(
	ctx context.Context,
	req *balancepb.DebitRequest,
) (*balancepb.DebitResponse, error) {
	res, err := h.svc.Debit(ctx, convert.DebitRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return convert.DebitResponseToProto(res), nil
}

func (h *Handler) Freeze(
	ctx context.Context,
	req *balancepb.FreezeRequest,
) (*balancepb.FreezeResponse, error) {
	res, err := h.svc.Freeze(ctx, convert.FreezeRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return convert.FreezeResponseToProto(res), nil
}

func (h *Handler) Unfreeze(
	ctx context.Context,
	req *balancepb.UnfreezeRequest,
) (*balancepb.UnfreezeResponse, error) {
	res, err := h.svc.Unfreeze(ctx, convert.UnfreezeRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return convert.UnfreezeResponseToProto(res), nil
}
