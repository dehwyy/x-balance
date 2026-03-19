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
		Id:            e.ID.Value,
		UserId:        e.UserID.Value,
		Type:          e.Type.Value,
		Amount:        e.Amount.Value.String(),
		TransactionId: e.TransactionID.Value,
		CreatedAt:     timestamppb.New(e.CreatedAt),
	}
}

func ProtoToEvent(p *transactionv1.Transaction) *event.Event {
	amount, _ := decimal.NewFromString(p.Amount)

	return &event.Event{
		ID:            event.ID{Value: p.Id},
		UserID:        user.ID{Value: p.UserId},
		Type:          event.EventType{Value: p.Type},
		Amount:        event.Amount{Value: amount},
		TransactionID: event.TransactionID{Value: p.TransactionId},
		CreatedAt:     p.CreatedAt.AsTime(),
	}
}
