package convert

import (
	"github.com/dehwyy/x-balance/internal/application/service/transactionservice"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	balancepb "github.com/dehwyy/x-balance/internal/generated/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ListTransactionsRequestToDomain(req *balancepb.ListTransactionsRequest) transactionservice.ListTransactionsRequest {
	r := transactionservice.ListTransactionsRequest{
		UserID: req.UserId,
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}
	if req.From != nil {
		t := req.From.AsTime()
		r.From = &t
	}
	if req.To != nil {
		t := req.To.AsTime()
		r.To = &t
	}
	return r
}

func ListTransactionsResponseToProto(res *transactionservice.ListTransactionsResponse) *balancepb.ListTransactionsResponse {
	txs := make([]*balancepb.Transaction, len(res.Events))
	for i, e := range res.Events {
		txs[i] = eventToProto(e)
	}
	return &balancepb.ListTransactionsResponse{
		Transactions: txs,
		Total:        res.Total,
	}
}

func GetTransactionRequestToDomain(req *balancepb.GetTransactionRequest) transactionservice.GetTransactionRequest {
	return transactionservice.GetTransactionRequest{
		UserID: req.UserId,
		TxID:   req.TxId,
	}
}

func GetTransactionResponseToProto(res *transactionservice.GetTransactionResponse) *balancepb.GetTransactionResponse {
	return &balancepb.GetTransactionResponse{
		Transaction: eventToProto(res.Event),
	}
}

func eventToProto(e *event.Event) *balancepb.Transaction {
	return &balancepb.Transaction{
		Id:            e.ID.Value,
		UserId:        e.UserID,
		Type:          e.Type.Value,
		Amount:        e.Amount.Value.String(),
		TransactionId: e.TransactionID.Value,
		CreatedAt:     timestamppb.New(e.CreatedAt),
	}
}
