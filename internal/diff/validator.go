package diff

import (
	"errors"
	"fmt"

	"github.com/acme/driftctl-diff/internal/drift"
)

// ValidationResult holds the outcome of validating a set of drift results.
type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

// Validator checks drift results for structural integrity and data quality.
type Validator struct {
	maxChangesPerResource int
}

// NewValidator creates a Validator with sensible defaults.
func NewValidator(maxChangesPerResource int) *Validator {
	if maxChangesPerResource <= 0 {
		maxChangesPerResource = 500
	}
	return &Validator{maxChangesPerResource: maxChangesPerResource}
}

// Validate inspects results and returns a ValidationResult.
func (v *Validator) Validate(results []drift.ResourceDiff) (ValidationResult, error) {
	if results == nil {
		return ValidationResult{}, errors.New("results must not be nil")
	}

	vr := ValidationResult{Valid: true}

	for i, r := range results {
		if r.ResourceID == "" {
			vr.Errors = append(vr.Errors, fmt.Sprintf("result[%d]: empty resource_id", i))
			vr.Valid = false
		}
		if r.ResourceType == "" {
			vr.Errors = append(vr.Errors, fmt.Sprintf("result[%d]: empty resource_type", i))
			vr.Valid = false
		}
		if len(r.Changes) > v.maxChangesPerResource {
			vr.Warnings = append(vr.Warnings,
				fmt.Sprintf("result[%d] (%s/%s): %d changes exceeds recommended max %d",
					i, r.ResourceType, r.ResourceID, len(r.Changes), v.maxChangesPerResource))
		}
	}

	return vr, nil
}
