package balancehandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	balanceconvert "github.com/dehwyy/x-balance/internal/delivery/api/balance/convert"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
)

func (h *Handler) Credit(
	ctx context.Context,
	req *balancev1.CreditRequest,
) (*balancev1.CreditResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceDelivery.Credit", dspan.Attr("req", req))
	defer span.End()

	response, err := h.balanceservice.Credit(ctx, balanceconvert.CreditRequestToDomain(req))
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := balanceconvert.CreditResponseToProto(response)
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
