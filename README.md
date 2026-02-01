# Infracost Hetzner ☁️💰

A lightweight CLI tool to estimate Hetzner Cloud costs directly from your Terraform plans. Designed for CI/CD pipelines (GitHub Actions) to prevent billing shocks before merging code.

![Go Version](https://img.shields.io/badge/go-1.21-blue)
![Docker](https://img.shields.io/badge/docker-ready-blue)

## 🚀 Key Features

*   **Real-time Pricing**: Fetches current prices directly from the Hetzner Cloud API.
*   **Robust Fallback**: Includes an embedded price table (offline mode) so your CI never fails due to API issues.
*   **Zero Dependencies**: Distributes as a static binary or Docker container.
*   **GitHub Integration**: Automatically posts a detailed cost table to your Pull Request.

## 📦 Installation

### Docker (Recommended)

```bash
docker pull fsz-codeshop/infracost-hetzner:latest
```

### From Source

Requirements: Go 1.21+

```bash
git clone https://github.com/fsz-codeshop/infracost-hetzner
cd infracost-hetzner
go build -o infracost-hetzner main.go
```

## 🛠️ Usage

### Local Testing

1.  Generate a Terraform plan JSON:
    ```bash
    cd terraform/
    terraform plan -out=plan.out
    terraform show -json plan.out > plan.json
    ```

2.  Run the tool:
    ```bash
    # Using the binary
    export HCLOUD_TOKEN="your_token"
    ./infracost-hetzner --plan plan.json

    # Using Docker
    docker run --rm -v $(pwd):/app \
      -e HCLOUD_TOKEN="your_token" \
      fsz-codeshop/infracost-hetzner --plan /app/plan.json
    ```

### 🤖 GitHub Actions Integration

Add this step to your `.github/workflows/pipeline.yml`:

```yaml
- name: Estimate Hetzner Costs
  uses: docker://fsz-codeshop/infracost-hetzner:latest
  env:
    HCLOUD_TOKEN: ${{ secrets.HCLOUD_TOKEN }}  # Required for real-time prices
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}  # Required for PR comments
    PR_NUMBER: ${{ github.event.pull_request.number }}
    GITHUB_REPOSITORY: ${{ github.repository }}
  with:
    args: --plan plan.json
```

*Note: Ensure you generate the `plan.json` in a previous step!*

## 🏗️ Architecture

1.  **Parser**: Reads standard `terraform show -json` output.
2.  **Engine**:
    *   Tries to fetch price from Hetzner API.
    *   If API fails/offline, loads embedded `fallback_prices.json`.
3.  **Reporter**: Formats output to Markdown and posts to GitHub.

## 🤝 Contributing

1.  Fork it
2.  Create your feature branch (`git checkout -b feature/my-new-feature`)
3.  Commit your changes (`git commit -am 'feat: add some feature'`)
4.  Push to the branch (`git push origin feature/my-new-feature`)
5.  Create new Pull Request

## 📄 License

MIT
