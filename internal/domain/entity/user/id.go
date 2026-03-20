package user

type ID struct {
	Value string
}

func NewID(v string) ID {
	return ID{
		Value: v,
	}
}
