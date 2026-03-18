package transactionhandler

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/service/transactionservice"
	"github.com/dehwyy/x-balance/internal/delivery/api/transaction/convert"
	balancepb "github.com/dehwyy/x-balance/internal/generated/pb"
)

type Handler struct {
	balancepb.UnimplementedTransactionServiceServer
	svc *transactionservice.Service
}

func New(svc *transactionservice.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListTransactions(
	ctx context.Context,
	req *balancepb.ListTransactionsRequest,
) (*balancepb.ListTransactionsResponse, error) {
	res, err := h.svc.ListTransactions(ctx, convert.ListTransactionsRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return convert.ListTransactionsResponseToProto(res), nil
}

func (h *Handler) GetTransaction(
	ctx context.Context,
	req *balancepb.GetTransactionRequest,
) (*balancepb.GetTransactionResponse, error) {
	res, err := h.svc.GetTransaction(ctx, convert.GetTransactionRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return convert.GetTransactionResponseToProto(res), nil
}
