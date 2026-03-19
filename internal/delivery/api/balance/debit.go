package balancehandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	balanceconvert "github.com/dehwyy/x-balance/internal/delivery/api/balance/convert"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
)

func (h *Handler) Debit(
	ctx context.Context,
	req *balancev1.DebitRequest,
) (*balancev1.DebitResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceDelivery.Debit", dspan.Attr("req", req))
	defer span.End()

	response, err := h.balanceservice.Debit(ctx, balanceconvert.DebitRequestToDomain(req))
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := balanceconvert.DebitResponseToProto(response)
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
