package event

type SnapshotID struct{ Value string }

func NewSnapshotID(v string) SnapshotID { return SnapshotID{Value: v} }
