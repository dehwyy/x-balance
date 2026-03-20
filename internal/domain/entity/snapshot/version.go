package snapshot

type Version struct {
	Value int64
}

func NewVersion(v int64) Version { return Version{Value: v} }
