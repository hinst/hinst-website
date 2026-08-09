---
name: generate-config-doc
description: Generate config documentation for Golang backend
---

# Generate backend config documentation

Gather information about environment variables referenced in the project. Search all `.go` files for these calls:
1. `gophers.ReadEnvVar`
1. `gophers.RequireEnvVar`
1. `os.Getenv`

Gather information about command line flags referenced in the project. Search for all `.go` files for `flag.` calls.

Some values are required only in certain modes. This information must also be checked and included.
For example: env var A is only needed if flag B is set.

Fill `# Config` section in `README.MD` explaining available command line flags and environment variables.
If `README.MD` lacks `# Config` section, then insert it.
If `README.MD` already contains `# Config` section, then replace it with the fresh documentation.
