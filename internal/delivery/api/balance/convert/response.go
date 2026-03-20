package balanceconvert

import (
	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
)

func CreditResponseToProto(resp *balanceservice.CreditResponse) *balancev1.CreditResponse {
	return &balancev1.CreditResponse{
		NewBalance:    resp.NewBalance.String(),
		TransactionId: string(resp.TransactionID),
	}
}

func DebitResponseToProto(resp *balanceservice.DebitResponse) *balancev1.DebitResponse {
	return &balancev1.DebitResponse{
		NewBalance:    resp.NewBalance.String(),
		TransactionId: string(resp.TransactionID),
	}
}

func FreezeResponseToProto(resp *balanceservice.FreezeResponse) *balancev1.FreezeResponse {
	return &balancev1.FreezeResponse{
		FrozenAmount:  resp.FrozenAmount.String(),
		TransactionId: string(resp.TransactionID),
	}
}

func UnfreezeResponseToProto(resp *balanceservice.UnfreezeResponse) *balancev1.UnfreezeResponse {
	return &balancev1.UnfreezeResponse{
		UnfrozenAmount: resp.UnfrozenAmount.String(),
		TransactionId:  string(resp.TransactionID),
	}
}

func GetBalanceResponseToProto(resp *balanceservice.GetBalanceResponse) *balancev1.GetBalanceResponse {
	return &balancev1.GetBalanceResponse{
		Available: resp.Available.String(),
		Frozen:    resp.Frozen.String(),
		Total:     resp.Total.String(),
	}
}
