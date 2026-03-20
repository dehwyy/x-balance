package event

type TransactionID string

func (t TransactionID) String() string {
	return string(t)
}

func (t TransactionID) ReleaseKey() TransactionID {
	return TransactionID(string(t) + ":release")
}
