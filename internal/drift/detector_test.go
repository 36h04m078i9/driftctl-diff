package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetect_NoDrift(t *testing.T) {
	detector := NewDetector()

	desired := map[string]interface{}{
		"instance_type": "t2.micro",
		"ami":           "ami-12345678",
	}
	actual := map[string]interface{}{
		"instance_type": "t2.micro",
		"ami":           "ami-12345678",
	}

	drifts := detector.Detect("aws_instance.web", desired, actual)
	assert.Empty(t, drifts)
}

func TestDetect_ValueChanged(t *testing.T) {
	detector := NewDetector()

	desired := map[string]interface{}{
		"instance_type": "t2.micro",
	}
	actual := map[string]interface{}{
		"instance_type": "t2.large",
	}

	drifts := detector.Detect("aws_instance.web", desired, actual)
	assert.Len(t, drifts, 1)
	assert.Equal(t, "aws_instance.web", drifts[0].ResourceID)
	assert.Equal(t, "instance_type", drifts[0].Attribute)
	assert.Equal(t, "t2.micro", drifts[0].Expected)
	assert.Equal(t, "t2.large", drifts[0].Actual)
}

func TestDetect_MissingAttribute(t *testing.T) {
	detector := NewDetector()

	desired := map[string]interface{}{
		"instance_type": "t2.micro",
		"tags":          map[string]interface{}{"Env": "prod"},
	}
	actual := map[string]interface{}{
		"instance_type": "t2.micro",
	}

	drifts := detector.Detect("aws_instance.web", desired, actual)
	assert.Len(t, drifts, 1)
	assert.Equal(t, "tags", drifts[0].Attribute)
	assert.Nil(t, drifts[0].Actual)
}

func TestDetect_MultipleChanges(t *testing.T) {
	detector := NewDetector()

	desired := map[string]interface{}{
		"instance_type": "t2.micro",
		"ami":           "ami-old",
	}
	actual := map[string]interface{}{
		"instance_type": "t3.medium",
		"ami":           "ami-new",
	}

	drifts := detector.Detect("aws_instance.web", desired, actual)
	assert.Len(t, drifts, 2)
}
