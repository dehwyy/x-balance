package transactionconvert

import (
	"github.com/dehwyy/x-balance/internal/application/service/transactionservice"
	eventconvert "github.com/dehwyy/x-balance/internal/domain/entity/event/convert"
	transactionv1 "github.com/dehwyy/x-balance/internal/generated/pb/common/transaction/v1"
	transactionspb "github.com/dehwyy/x-balance/internal/generated/pb/transactions/v1"
)

func ListTransactionsResponseToProto(resp *transactionservice.ListTransactionsResponse) *transactionspb.ListTransactionsResponse {
	txs := make([]*transactionv1.Transaction, len(resp.Events))
	for i, e := range resp.Events {
		txs[i] = eventconvert.EventToProto(e)
	}
	return &transactionspb.ListTransactionsResponse{
		Transactions: txs,
		Total:        resp.Total,
	}
}

func GetTransactionResponseToProto(resp *transactionservice.GetTransactionResponse) *transactionspb.GetTransactionResponse {
	return &transactionspb.GetTransactionResponse{
		Transaction: eventconvert.EventToProto(resp.Event),
	}
}
