package transactionhandler

import (
	"go.uber.org/fx"

	"github.com/dehwyy/x-balance/internal/application/service/transactionservice"
	transactionspb "github.com/dehwyy/x-balance/internal/generated/pb/transactions/v1"
)

type Opts struct {
	fx.In
	TransactionService *transactionservice.Service
}

type Handler struct {
	transactionspb.UnimplementedTransactionServiceServer
	transactionservice *transactionservice.Service
}

func New(opts Opts) *Handler {
	return &Handler{
		transactionservice: opts.TransactionService,
	}
}
