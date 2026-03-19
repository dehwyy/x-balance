package balancehandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	balanceconvert "github.com/dehwyy/x-balance/internal/delivery/api/balance/convert"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
)

func (h *Handler) GetBalance(
	ctx context.Context,
	req *balancev1.GetBalanceRequest,
) (*balancev1.GetBalanceResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceDelivery.GetBalance", dspan.Attr("req", req))
	defer span.End()

	response, err := h.balanceservice.GetBalance(ctx, balanceconvert.GetBalanceRequestToDomain(req))
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := balanceconvert.GetBalanceResponseToProto(response)
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
