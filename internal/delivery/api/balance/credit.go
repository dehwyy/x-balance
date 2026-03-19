package balancehandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
	"github.com/shopspring/decimal"
)

func (h *Handler) Credit(ctx context.Context, req *balancev1.CreditRequest) (*balancev1.CreditResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceDelivery.Credit", dspan.Attr("req", req))
	defer span.End()

	amount, _ := decimal.NewFromString(req.Amount)

	response, err := h.balanceservice.Credit(ctx, balanceservice.CreditRequest{
		UserID:        req.UserId,
		Amount:        amount,
		TransactionID: req.TransactionId,
	})
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := &balancev1.CreditResponse{
		NewBalance:    response.NewBalance.String(),
		TransactionId: response.TransactionID,
	}
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
