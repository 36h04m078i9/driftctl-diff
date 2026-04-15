// Package drift defines shared types used across drift detection.
package drift

// Change describes a single attribute-level difference between Terraform
// state and the live cloud resource.
type Change struct {
	// Attribute is the resource attribute key that differs.
	Attribute string
	// StateValue is the value recorded in Terraform state.
	StateValue string
	// LiveValue is the value observed from the live cloud resource.
	LiveValue string
}

// ResourceDiff groups all Changes for a single resource.
type ResourceDiff struct {
	// ResourceType is the Terraform resource type, e.g. "aws_instance".
	ResourceType string
	// ResourceID is the unique identifier of the resource.
	ResourceID string
	// Changes is the list of attribute-level differences. An empty slice
	// means the resource is in sync.
	Changes []Change
}
