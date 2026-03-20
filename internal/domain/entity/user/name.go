package user

type Name struct {
	Value string
}

func NewName(v string) Name {
	return Name{
		Value: v,
	}
}
