package balancehandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
)

func (h *Handler) GetBalance(ctx context.Context, req *balancev1.GetBalanceRequest) (*balancev1.GetBalanceResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceDelivery.GetBalance", dspan.Attr("req", req))
	defer span.End()

	response, err := h.balanceservice.GetBalance(ctx, balanceservice.GetBalanceRequest{UserID: req.UserId})
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := &balancev1.GetBalanceResponse{
		Available: response.Available.String(),
		Frozen:    response.Frozen.String(),
		Total:     response.Total.String(),
	}
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
