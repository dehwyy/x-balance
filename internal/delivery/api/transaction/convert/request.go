package transactionconvert

import (
	"github.com/dehwyy/x-balance/internal/application/service/transactionservice"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	transactionspb "github.com/dehwyy/x-balance/internal/generated/pb/transactions/v1"
	"github.com/dehwyy/x-balance/pkg/storage"
)

func ListTransactionsRequestToDomain(req *transactionspb.ListTransactionsRequest) *transactionservice.ListTransactionsRequest {
	r := &transactionservice.ListTransactionsRequest{
		UserID:     user.NewID(req.UserId),
		Pagination: storage.NewPagination(int(req.Limit), int(req.Offset)),
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

func GetTransactionRequestToDomain(req *transactionspb.GetTransactionRequest) *transactionservice.GetTransactionRequest {
	return &transactionservice.GetTransactionRequest{
		UserID: user.NewID(req.UserId),
		TxID:   event.NewID(req.TxId),
	}
}
