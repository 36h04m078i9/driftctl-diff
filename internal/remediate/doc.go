// Package remediate provides remediation suggestions for infrastructure drift
// detected between Terraform state and live cloud resources.
//
// Usage:
//
//	suggestions := remediate.Suggest(results)
//	p := remediate.New(os.Stdout)
//	p.Print(suggestions)
//
// Each Suggestion describes which resource is drifted and provides a
// human-readable hint pointing operators toward `terraform apply`.
package remediate
