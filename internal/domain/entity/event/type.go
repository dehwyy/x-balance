package event

type EventType struct {
	Value string
}

var (
	TypeCredit        = EventType{Value: "credit"}
	TypeDebit         = EventType{Value: "debit"}
	TypeFreezeHold    = EventType{Value: "freeze_hold"}
	TypeFreezeRelease = EventType{Value: "freeze_release"}
)
