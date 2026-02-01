# Technical Specifications - Infracost Hetzner

## 1. Architecture Overview
The tool is a stateless CLI written in Go. It follows a simple Pipeline pattern:
`Input (Terraform JSON) -> Parser (Go) -> Pricing Engine (API/Hcloud-Go) -> Formatter (Markdown/Table) -> Output (Stdout/GitHub API)`

## 2. Tech Stack
- **Language**: Go 1.21+
- **Primary Libs**: 
    - `github.com/hetznercloud/hcloud-go/v2`: Official Hetzner SDK.
    - `github.com/google/go-github/v60`: For PR comments.
    - `github.com/spf13/cobra`: CLI framework.
- **Packaging**: Docker (multi-stage build with Alpine or Scratch).

## 3. Implementation Plan

### Phase 1: Foundation (The Parser)
- Initialize Go module.
- Create struct mappings for `terraform show -json` output.
- Focus on extracting `resource_changes` where `type` starts with `hcloud_`.

### Phase 2: Pricing Engine
- Implement a `PriceProvider` interface.
- **HcloudAPIProvider**: Real-time fetch using `hcloud-go`.
- **LocalFallbackProvider**: Uses `embed` feature of Go (1.16+) to bake a `prices.json` into the binary.
- Logic: `if err := api.Fetch(); err != nil { useFallback() }`.

### Phase 3: Resource Mapping
- Implement cost calculation for:
    - **Servers**: Based on `server_type` and `location`/`datacenter`.
    - **Volumes**: Based on size (GB).
    - **Load Balancers**: Based on type.
- Future-proofing for **Floating IPs** and **Primary IPs**.

### Phase 4: Integration
- Configure GitHub Actions environment variables (`GITHUB_TOKEN`, `HCLOUD_TOKEN`).
- Implement the "GitHub Commenter" logic using the PR context (Repo, Owner, PR Number) usually provided by CI environment variables.

## 4. Security & Performance
- **Secrets**: `HCLOUD_TOKEN` must be handled as a sensitive secret.
- **Binary Size**: Use `-ldflags="-s -w"` to reduce Docker layer size.
- **Concurrency**: Use Goroutines for fetching prices of multiple resource types if needed, though sequential should be fast enough initially.

## 5. Data Schema (Internal Fallback)
```json
{
  "version": "2024-02-01",
  "servers": {
    "cx11": { "monthly": 4.15, "hourly": 0.007 },
    "cpx11": { "monthly": 4.75, "hourly": 0.008 }
  }
}
```
