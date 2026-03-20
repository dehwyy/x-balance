package eventconvert

import (
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	transactionv1 "github.com/dehwyy/x-balance/internal/generated/pb/common/transaction/v1"
)

func EventToProto(e *event.Event) *transactionv1.Transaction {
	return &transactionv1.Transaction{
		Id:            string(e.ID),
		UserId:        string(e.UserID),
		Type:          e.Type,
		Amount:        decimal.Decimal(e.Amount).String(),
		TransactionId: string(e.TransactionID),
		CreatedAt:     timestamppb.New(e.CreatedAt),
	}
}

func ProtoToEvent(p *transactionv1.Transaction) *event.Event {
	amount, _ := decimal.NewFromString(p.Amount)

	return &event.Event{
		ID:            event.ID(p.Id),
		UserID:        user.ID(p.UserId),
		Type:          p.Type,
		Amount:        event.Amount(amount),
		TransactionID: event.TransactionID(p.TransactionId),
		CreatedAt:     p.CreatedAt.AsTime(),
	}
}
