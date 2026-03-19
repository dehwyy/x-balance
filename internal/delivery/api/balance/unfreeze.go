package balancehandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	balanceconvert "github.com/dehwyy/x-balance/internal/delivery/api/balance/convert"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
)

func (h *Handler) Unfreeze(
	ctx context.Context,
	req *balancev1.UnfreezeRequest,
) (*balancev1.UnfreezeResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceDelivery.Unfreeze", dspan.Attr("req", req))
	defer span.End()

	response, err := h.balanceservice.Unfreeze(ctx, balanceconvert.UnfreezeRequestToDomain(req))
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := balanceconvert.UnfreezeResponseToProto(response)
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
