// Package config provides loading and validation of driftctl-diff configuration.
//
// Configuration can be supplied via a YAML file (default: .driftctl.yaml in the
// working directory) or overridden through CLI flags.  The [Load] function
// merges a file on top of [DefaultConfig], so every field always has a
// well-defined value even when the file omits it.
//
// Example YAML:
//
//	state_path: path/to/terraform.tfstate
//	provider: aws
//	region: us-east-1
//	output_format: text   # text | json
//	color: true
//	ignore:
//	  - aws_instance.bastion
package config
