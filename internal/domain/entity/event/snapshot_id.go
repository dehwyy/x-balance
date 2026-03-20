package event

type SnapshotID string

func (s SnapshotID) String() string {
	return string(s)
}
