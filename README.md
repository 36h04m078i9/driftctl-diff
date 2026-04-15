# driftctl-diff

> A CLI tool that surfaces infrastructure drift between Terraform state and live cloud resources in a human-readable diff format.

---

## Installation

**Using Go:**
```bash
go install github.com/yourusername/driftctl-diff@latest
```

**Using Homebrew:**
```bash
brew install yourusername/tap/driftctl-diff
```

---

## Usage

Run against your current Terraform working directory:

```bash
driftctl-diff --provider aws --region us-east-1
```

**Example output:**

```diff
~ aws_security_group.web (sg-0abc123)
  - ingress.0.cidr_blocks: ["10.0.0.0/8"]
  + ingress.0.cidr_blocks: ["0.0.0.0/0"]

+ aws_s3_bucket.untracked-bucket (not in state)

- aws_iam_role.deprecated (in state, missing in cloud)
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--provider` | Cloud provider (`aws`, `gcp`, `azure`) | `aws` |
| `--region` | Target region | `us-east-1` |
| `--state` | Path to Terraform state file | `terraform.tfstate` |
| `--output` | Output format (`diff`, `json`, `summary`) | `diff` |

---

## Requirements

- Go 1.21+
- Terraform state file or remote backend access
- Appropriate cloud provider credentials configured

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE)