# PRD - Infracost Hetzner (Go Edition)

## 1. Executive Summary
**Problem Statement**: Platform engineers and developers lack immediate visibility into the financial impact of infrastructure changes on Hetzner Cloud during the Pull Request process.

**Proposed Solution**: A lightweight Go-based Dockerized CLI tool that parses Terraform plans, fetches real-time prices from Hetzner API (with local fallback), and comments the estimated cost directly on GitHub PRs.

**Business Impact**: 
- Prevent "billing shocks" before infrastructure is provisioned.
- Increase financial awareness among developers.
- Automate cost governance in the CI/CD pipeline.

## 2. User Stories
- **As a Developer**, I want to see how much my PR will increase the monthly bill so I can optimize my resource choices.
- **As a DevOps Engineer**, I want the CI to fail or warn me if the cost exceeds a certain threshold (future) or if price fetching fails.
- **As a FinOps Lead**, I want to ensure that prices being report are accurate or at least clearly marked as "estimated fallback".

## 3. Scope
### 3.1 In Scope
- **Input**: Terraform plan in JSON format.
- **Processing**: Mapping standard Hetzner Terraform resources (`hcloud_server`, `hcloud_volume`, `hcloud_load_balancer`) to their respective costs.
- **Price Engine**: API first (via `hcloud-go`), Fallback to embedded JSON/CSV for disconnected/error states.
- **Output**: 
    - Detailed CLI Table (Console).
    - GitHub PR Comment (Markdown Table).
- **Core Resources**: Servers (Types/Locations), Volumes, Load Balancers.

### 3.2 Out of Scope
- Support for other cloud providers (Hetzner only).
- Direct budget enforcement (initial version).
- Historical cost analysis (focus is on "change" cost).

## 4. Success Metrics
- **Performance**: Execution time < 10 seconds (excluding Terraform plan generation).
- **Reliability**: 100% success rate on PR commenting if token is valid.
- **Accuracy**: Cost estimation matching Hetzner's public pricing within 5% margin (considering rounding).

## 5. Decision Log
- **Language**: Go. Chosen for performance, static binaries, and official SDK support.
- **Fallback Mechanism**: Integrated embedded data to ensure pipeline stability.
- **Platform**: Docker + GitHub Actions. Targets the most common modern developer workflow.
