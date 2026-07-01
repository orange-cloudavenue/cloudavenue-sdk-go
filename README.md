# Cloud Avenue SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/orange-cloudavenue/cloudavenue-sdk-go.svg)](https://pkg.go.dev/github.com/orange-cloudavenue/cloudavenue-sdk-go)

> **⚠️ Note:** This SDK is currently undergoing a major rewrite. A new version is coming soon — but until then, this is the stable release to use.

[![Go version](https://img.shields.io/github/go-mod/go-version/orange-cloudavenue/cloudavenue-sdk-go)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/orange-cloudavenue/cloudavenue-sdk-go)](https://github.com/orange-cloudavenue/cloudavenue-sdk-go/releases)
[![License: MPL-2.0](https://img.shields.io/badge/License-MPL--2.0-blue.svg)](https://opensource.org/licenses/MPL-2.0)
[![CI](https://github.com/orange-cloudavenue/cloudavenue-sdk-go/actions/workflows/pr.yml/badge.svg)](https://github.com/orange-cloudavenue/cloudavenue-sdk-go/actions/workflows/pr.yml)

**cloudavenue-sdk-go** is the official Go SDK for the [**Cloud Avenue**](https://www.orange-business.com/en/our-solutions/cloud/cloud-avenue) platform by Orange Business. It provides a Go-native client library to programmatically manage your Cloud Avenue infrastructure — Edge Gateways, VDCs, VDC Groups, networking, security, load balancing, S3-compatible object storage, backup, and more.

The SDK wraps and extends the upstream [VMware go-vcloud-director SDK](https://github.com/vmware/go-vcloud-director) with [Cloud Avenue-specific APIs (InfrAPI)](https://www.orange-business.com/en/solutions/apiforbusiness/cloud-avenue-api) for resource types not available through standard VMware VCD endpoints.

---

## Features

- **Full lifecycle management** for Edge Gateways, VDCs, VDC Groups, Tier-0 VRFs, public IPs, and more
- **Advanced networking** — firewall rules, security groups, IPSets, application port profiles, network context profiles, isolated/routed networks
- **Load balancing** — Advanced Load Balancer (ALB) pools, virtual services, HTTP request/response/security policies
- **S3-compatible object storage** — buckets, credentials, users via the AWS S3 API
- **Backup management** — Veritas NetBackup integration for inventory, machines, and protection levels
- **IAM** — local and SAML user management (CRUD, password, enable/disable)
- **Bare Metal Servers** — inventory and hostname-based lookup
- **VCDA (DRaaS)** — IP allowlisting for VMware Cloud Director Availability
- **Job-based async operations** — with context-aware wait, configurable polling, and timeout support
- **Concurrent dual-backend lookups** — queries VMware govcd and Cloud Avenue InfrAPI in parallel for optimal latency
- **Regional console routing** — automatic organization-to-console mapping with per-console service availability

---

## Installation

```bash
go get github.com/orange-cloudavenue/cloudavenue-sdk-go@latest
```

Requires Go **1.23** or later.

---

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/orange-cloudavenue/cloudavenue-sdk-go"
)

func main() {
    // Create a client with environment variable configuration.
    // Required env vars:
    //   CLOUDAVENUE_USERNAME
    //   CLOUDAVENUE_PASSWORD
    //   CLOUDAVENUE_ORG
    //   CLOUDAVENUE_URL (e.g. https://console1.cloudavenue.orange-business.com)
    client, err := cloudavenue.New()
    if err != nil {
        log.Fatalf("failed to create client: %v", err)
    }

    // List all Edge Gateways
    edgeGateways, err := client.V1().EdgeGateway.List(context.Background())
    if err != nil {
        log.Fatalf("failed to list edge gateways: %v", err)
    }
    for _, egw := range edgeGateways {
        fmt.Printf("Edge Gateway: %s (%s)\n", egw.EdgeGateway.Name, egw.EdgeGateway.ID)
    }

    // Get a VDC and inspect its properties
    vdc, err := client.V1().VDC().Get(context.Background(), "my-vdc")
    if err != nil {
        log.Fatalf("failed to get VDC: %v", err)
    }
    fmt.Printf("VDC: %s, CPU: %d MHz, Memory: %d MB\n",
        vdc.VDC.Vdc.Name, vdc.GetVCPUInMhz(), vdc.GetMemory())
}
```

---

## Authentication

Configure via **environment variables**:

| Variable               | Description                           | Required |
| ---------------------- | ------------------------------------- | -------- |
| `CLOUDAVENUE_URL`      | Cloud Avenue console URL              | Yes      |
| `CLOUDAVENUE_USERNAME` | API username                          | Yes      |
| `CLOUDAVENUE_PASSWORD` | API password                          | Yes      |
| `CLOUDAVENUE_ORG`      | Organization name                     | Yes      |
| `CLOUDAVENUE_DEBUG`    | Enable debug logging (`true`/`false`) | No       |
| `CLOUDAVENUE_DEV`      | Development mode flag                 | No       |
| `CLOUDAVENUE_CORE_API` | Override backend API endpoint         | No       |

The SDK uses OAuth2 token-based authentication (v2, introduced in v0.27.0). Tokens are obtained automatically on client creation and refreshed as needed.

> **Note**: The legacy authentication method reached end of life on October 1, 2026. Please upgrade to a current SDK version.

---

## Client Architecture

```
Client (cloudavenue package)
├── V1()          → Top-level API surface
│   ├── VDC()             → VDC operations
│   ├── AdminVDC()        → Admin VDC operations
│   ├── EdgeGateway       → Edge Gateway CRUD, firewall, groups
│   ├── EdgeGateway.ALB   → Load balancer pools & virtual services
│   ├── T0                → Tier-0 VRF gateway bandwidth
│   ├── PublicIP          → Public IP management
│   ├── VCDA              → DRaaS IP allowlisting
│   ├── BMS               → Bare Metal Server inventory
│   ├── IAM()             → User management
│   ├── Org()             → Org properties & certificates
│   ├── AdminOrg()        → Admin org (catalogs, vApp leases)
│   ├── Querier()         → VMware query service
│   ├── Vmware()          → Direct VMware VCD access
│   ├── S3()              → S3 client
│   └── Netbackup         → NetBackup management
├── Config()      → Client configuration
└── S3()          → Standalone S3 client
```

---

## API Overview

### Infrastructure & Networking

| Package                | Resources                                                                                                          |
| ---------------------- | ------------------------------------------------------------------------------------------------------------------ |
| `v1/edgegw`            | Edge Gateway CRUD, bandwidth, firewall rules, security groups, IPSets, app port profiles, network context profiles |
| `v1/edgegateway/`      | Edge Gateway CRUD (refactored, interface-based), network services                                                  |
| `v1/edgeloadbalancer/` | ALB pools, virtual services, HTTP request/response/security policies                                               |
| `v1/vdc`               | VDC CRUD, storage profiles, security groups, IPSets, vApps, isolated/routed networks                               |
| `v1/vdcg`              | VDC Group management, distributed firewall, security groups, IPSets, networks, network context profiles            |
| `v1/t0`                | Tier-0 VRF gateway listing, bandwidth capacity, service classes                                                    |
| `v1/publicip`          | Public IP listing, creation, deletion, job tracking                                                                |
| `v1/vcda`              | VCDA IP allowlisting for DRaaS                                                                                     |
| `v1/bms`               | Bare Metal Server inventory                                                                                        |

### Security & IAM

| Package   | Resources                                                                                  |
| --------- | ------------------------------------------------------------------------------------------ |
| `v1/iam/` | Local & SAML users: create, read, update, delete, enable, disable, unlock, change password |

### Organization

| Package           | Resources                                                                          |
| ----------------- | ---------------------------------------------------------------------------------- |
| `v1/org/`         | Org properties (name, description, email, billing model), certificate library CRUD |
| `v1/admin_org.go` | Catalog listing, vApp lease settings                                               |

### Storage

| Package               | Resources                            |
| --------------------- | ------------------------------------ |
| `v1/s3_bucket.go`     | S3 bucket operations                 |
| `v1/s3_credential.go` | S3 credential management             |
| `v1/s3_user.go`       | S3 user listing, canonical ID lookup |

### Backup

| Package         | Resources                                                                          |
| --------------- | ---------------------------------------------------------------------------------- |
| `v1/netbackup/` | NetBackup inventory, machines, protect jobs, protection levels, vCloud integration |

### Query Service

| Package         | Resources                                                       |
| --------------- | --------------------------------------------------------------- |
| `v1/querier.go` | VMware query service — list/get VDCs, vApps, VMs, Edge Gateways |

---

## Usage Examples

### Edge Gateway: Manage Bandwidth

```go
// Get bandwidth capacity remaining
remaining, err := client.V1().EdgeGateway.GetBandwidthCapacityRemaining(context.Background(), "my-edgegw")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Remaining bandwidth: %d Mbps\n", remaining)

// Update bandwidth
err = client.V1().EdgeGateway.UpdateBandwidth(context.Background(), "my-edgegw", 500)
if err != nil {
    log.Fatal(err)
}
```

### Edge Gateway: Firewall Rules with Network Context Profiles

```go
// Use the Extended firewall methods to access networkContextProfiles
egw, err := client.V1().EdgeGateway.Get(context.Background(), "my-edgegw")
if err != nil {
    log.Fatal(err)
}

fw, err := egw.GetFirewallExtended(context.Background())
if err != nil {
    log.Fatal(err)
}

// rules now support NetworkContextProfiles
for _, rule := range fw.Rules {
    fmt.Printf("Rule %s: %s → %s (profiles: %v)\n",
        rule.Name, rule.Source, rule.Destination, rule.NetworkContextProfiles)
}
```

### VDC Group: Distributed Firewall

```go
vdcg, err := client.V1().VDC().GetVDCGroup(context.Background(), "my-vdcg")
if err != nil {
    log.Fatal(err)
}

fw, err := vdcg.GetFirewall(context.Background())
if err != nil {
    log.Fatal(err)
}

rules, err := fw.GetRules(context.Background())
if err != nil {
    log.Fatal(err)
}

for _, rule := range rules {
    fmt.Printf("Rule: %s | Action: %s | Source: %s → Dest: %s\n",
        rule.Name, rule.Action, rule.Source, rule.Destination)
}
```

### Asynchronous Jobs with Context

```go
// Operations like VDC/EdgeGateway creation return a JobStatus
job, err := client.V1().EdgeGateway.New(context.Background(), ...)
if err != nil {
    log.Fatal(err)
}

// Wait with a 2-minute timeout
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()

err = job.WaitWithContext(ctx)
if err != nil {
    log.Fatal(err)
}
```

### S3 Operations

```go
s3, err := client.V1().S3()
if err != nil {
    log.Fatal(err) // S3 may not be available in all regions
}

buckets, err := s3.ListBuckets(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, b := range buckets {
    fmt.Printf("Bucket: %s\n", b.Name)
}
```

---

## Development

### Prerequisites

- Go 1.23+
- [golangci-lint](https://golangci-lint.run/) (for linting)
- [pre-commit](https://pre-commit.com/) (optional, for commit hooks)
- [Task](https://taskfile.dev/) v3 (optional, for extended task runner)

### Commands

```bash
make lint        # Run golangci-lint
make fmt         # Format code with gofmt
make test        # Run all tests with coverage
make pre-commit  # Install pre-commit hooks

# Or using Task
task install     # Install tools and setup
task lint        # Run all linters
task generate    # Run code generation
```

### Testing

Tests use `httpmock` for HTTP mocking and `go.uber.org/mock` for generated interface mocks. Many packages provide `NewFakeClient` constructors for unit testing without real API credentials.

```bash
make test
# or
go test -v -cover -timeout=120s -parallel=4 ./...
```

### Code Generation

The project uses `go:generate` directives with `go.uber.org/mock` for mock generation. Run:

```bash
go generate ./...
```

---

## Project Structure

```
cloudavenue-sdk-go/
├── cloudavenue.go            # Root client entry point
├── v1/                       # Main API surface
│   ├── edgegw/               # Edge Gateway management
│   ├── edgegateway/          # Edge Gateway CRUD (refactored)
│   ├── edgeloadbalancer/     # ALB management
│   ├── iam/                  # IAM user management
│   ├── infrapi/              # Cloud Avenue InfrAPI client
│   ├── netbackup/            # NetBackup management
│   ├── org/                  # Organization operations
│   └── ...                   # Flat files for additional resources
├── pkg/                      # Shared utilities
│   ├── clients/              # Auth clients (cloudavenue, s3, netbackup, consoles)
│   ├── common/               # Shared types (API errors, job status)
│   ├── errors/               # Sentinel errors
│   ├── urn/                  # URN parsing, validation, normalization
│   └── helpers/              # Firewall and VDC group helpers
├── internal/                 # Internal implementation
│   ├── endpoints/            # API endpoint constants
│   └── utils/                # Generic helpers (ToPTR)
├── .github/                  # CI/CD workflows, dependabot, codeowners
└── .changelog/               # Individual changelog entries
```

---

## Key Design Patterns

- **Interface-based**: Core operations define `Client` interfaces with separate `goVCD` and `cloudavenue` sub-interfaces, enabling comprehensive mock generation for testing.
- **Dual-backend**: Resources are often fetched from both VMware govcd (VCD-native data) and Cloud Avenue InfrAPI (platform-specific properties) concurrently via `errgroup`.
- **Job-based async**: Long-running operations return `JobStatus` objects. Call `Wait()` or `WaitWithContext()` with configurable polling intervals and timeouts.
- **URN system**: A dedicated `pkg/urn` package validates and normalizes URNs for 20+ resource types. Used pervasively across the SDK.
- **Console routing**: Organizations are automatically mapped to regional consoles (Console1–Console9) via regex patterns. Each console tracks available services (S3, NetBackup, VCDA).

---

## Dependencies

| Dependency                                                                 | Purpose                 |
| -------------------------------------------------------------------------- | ----------------------- |
| [go-vcloud-director](https://github.com/vmware/go-vcloud-director) v2.26.1 | VMware VCD SDK          |
| [go-resty](https://github.com/go-resty/resty) v2.16.5                      | HTTP client for InfrAPI |
| [aws-sdk-go](https://github.com/aws/aws-sdk-go) v1.55.7                    | S3-compatible storage   |
| [envconfig](https://github.com/sethvargo/go-envconfig) v1.3.0              | Environment config      |

---

## Contributing

1. Ensure commits follow [Conventional Commits](https://www.conventionalcommits.org/) (enforced by pre-commit hooks)
2. Add changelog entries to `.changelog/` for each change
3. Run `make lint && make test` before opening a PR
4. All PRs require CI checks to pass (lint, unit tests, license check)

---

## License

Mozilla Public License 2.0. See [LICENSE](./LICENSE).

Copyright © Orange Business, 2025–2026.
