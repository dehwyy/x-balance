package balanceconvert

import (
	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
	"github.com/shopspring/decimal"
)

func CreditRequestToDomain(req *balancev1.CreditRequest) *balanceservice.CreditRequest {
	amount, _ := decimal.NewFromString(req.Amount)
	return &balanceservice.CreditRequest{
		UserID:        user.ID(req.UserId),
		Amount:        amount,
		TransactionID: event.TransactionID(req.TransactionId),
	}
}

func DebitRequestToDomain(req *balancev1.DebitRequest) *balanceservice.DebitRequest {
	amount, _ := decimal.NewFromString(req.Amount)
	return &balanceservice.DebitRequest{
		UserID:        user.ID(req.UserId),
		Amount:        amount,
		TransactionID: event.TransactionID(req.TransactionId),
	}
}

func FreezeRequestToDomain(req *balancev1.FreezeRequest) *balanceservice.FreezeRequest {
	amount, _ := decimal.NewFromString(req.Amount)
	return &balanceservice.FreezeRequest{
		UserID:               user.ID(req.UserId),
		Amount:               amount,
		TransactionID:        event.TransactionID(req.TransactionId),
		FreezeTimeoutSeconds: req.FreezeTimeoutSeconds,
	}
}

func UnfreezeRequestToDomain(req *balancev1.UnfreezeRequest) *balanceservice.UnfreezeRequest {
	return &balanceservice.UnfreezeRequest{
		UserID:        user.ID(req.UserId),
		TransactionID: event.TransactionID(req.TransactionId),
	}
}

func GetBalanceRequestToDomain(req *balancev1.GetBalanceRequest) *balanceservice.GetBalanceRequest {
	return &balanceservice.GetBalanceRequest{
		UserID: user.ID(req.UserId),
	}
}
