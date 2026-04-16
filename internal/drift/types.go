package drift

// ChangeKind describes the nature of an attribute-level drift.
type ChangeKind int

const (
	// KindChanged means the attribute exists in both state and live but values differ.
	KindChanged ChangeKind = iota
	// KindAdded means the attribute is present live but absent in state.
	KindAdded
	// KindDeleted means the attribute is present in state but absent live.
	KindDeleted
)

// String returns a human-readable label for a ChangeKind.
func (k ChangeKind) String() string {
	switch k {
	case KindChanged:
		return "changed"
	case KindAdded:
		return "added"
	case KindDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// AttributeChange captures a single attribute-level drift.
type AttributeChange struct {
	Attribute  string
	Kind       ChangeKind
	StateValue string
	LiveValue  string
}

// ResourceDiff groups all attribute changes for one resource.
type ResourceDiff struct {
	ResourceType string
	ResourceID   string
	Changes      []AttributeChange
}
