package convert

import (
	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	balancepb "github.com/dehwyy/x-balance/internal/generated/pb"
	"github.com/shopspring/decimal"
)

func GetBalanceResponseToProto(res *balanceservice.GetBalanceResponse) *balancepb.GetBalanceResponse {
	return &balancepb.GetBalanceResponse{
		Available: res.Available.String(),
		Frozen:    res.Frozen.String(),
		Total:     res.Total.String(),
	}
}

func CreditRequestToDomain(req *balancepb.CreditRequest) balanceservice.CreditRequest {
	amount, _ := decimal.NewFromString(req.Amount)
	return balanceservice.CreditRequest{
		UserID:        req.UserId,
		Amount:        amount,
		TransactionID: req.TransactionId,
	}
}

func CreditResponseToProto(res *balanceservice.CreditResponse) *balancepb.CreditResponse {
	return &balancepb.CreditResponse{
		NewBalance:    res.NewBalance.String(),
		TransactionId: res.TransactionID,
	}
}

func DebitRequestToDomain(req *balancepb.DebitRequest) balanceservice.DebitRequest {
	amount, _ := decimal.NewFromString(req.Amount)
	return balanceservice.DebitRequest{
		UserID:        req.UserId,
		Amount:        amount,
		TransactionID: req.TransactionId,
	}
}

func DebitResponseToProto(res *balanceservice.DebitResponse) *balancepb.DebitResponse {
	return &balancepb.DebitResponse{
		NewBalance:    res.NewBalance.String(),
		TransactionId: res.TransactionID,
	}
}

func FreezeRequestToDomain(req *balancepb.FreezeRequest) balanceservice.FreezeRequest {
	amount, _ := decimal.NewFromString(req.Amount)
	return balanceservice.FreezeRequest{
		UserID:               req.UserId,
		Amount:               amount,
		TransactionID:        req.TransactionId,
		FreezeTimeoutSeconds: req.FreezeTimeoutSeconds,
	}
}

func FreezeResponseToProto(res *balanceservice.FreezeResponse) *balancepb.FreezeResponse {
	return &balancepb.FreezeResponse{
		FrozenAmount:  res.FrozenAmount.String(),
		TransactionId: res.TransactionID,
	}
}

func UnfreezeRequestToDomain(req *balancepb.UnfreezeRequest) balanceservice.UnfreezeRequest {
	return balanceservice.UnfreezeRequest{
		UserID:        req.UserId,
		TransactionID: req.TransactionId,
	}
}

func UnfreezeResponseToProto(res *balanceservice.UnfreezeResponse) *balancepb.UnfreezeResponse {
	return &balancepb.UnfreezeResponse{
		UnfrozenAmount: res.UnfrozenAmount.String(),
		TransactionId:  res.TransactionID,
	}
}
