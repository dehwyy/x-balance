package balancehandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	balanceconvert "github.com/dehwyy/x-balance/internal/delivery/api/balance/convert"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
)

func (h *Handler) Freeze(
	ctx context.Context,
	req *balancev1.FreezeRequest,
) (*balancev1.FreezeResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceDelivery.Freeze", dspan.Attr("req", req))
	defer span.End()

	response, err := h.balanceservice.Freeze(ctx, balanceconvert.FreezeRequestToDomain(req))
	if err != nil {
		return nil, span.Err(err)
	}

	return dspan.Response(span, balanceconvert.FreezeResponseToProto(response)), nil
}
