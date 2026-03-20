package event

type TransactionID struct {
	Value string
}

func NewTransactionID(v string) TransactionID { return TransactionID{Value: v} }

func (t TransactionID) ReleaseKey() TransactionID { return TransactionID{Value: t.Value + ":release"} }
