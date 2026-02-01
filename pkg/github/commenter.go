package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/fsz-codeshop/infracost-hetzner/pkg/config"
	"github.com/fsz-codeshop/infracost-hetzner/pkg/pricing"
	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

func CommentPR(total *pricing.TotalCost, cfg *config.Config) error {
	if cfg.GithubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN not set")
	}

	if cfg.GithubRepo == "" {
		return fmt.Errorf("GITHUB_REPOSITORY not set")
	}

	parts := strings.Split(cfg.GithubRepo, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid GITHUB_REPOSITORY format")
	}
	owner, repo := parts[0], parts[1]

	if cfg.PRNumber == "" {
		return fmt.Errorf("PR_NUMBER not set")
	}

	prNumber, err := strconv.Atoi(cfg.PRNumber)
	if err != nil {
		return fmt.Errorf("invalid PR number: %s", cfg.PRNumber)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	body := formatMarkdown(total)
	comment := &github.IssueComment{
		Body: github.String(body),
	}

	_, _, err = client.Issues.CreateComment(ctx, owner, repo, prNumber, comment)
	return err
}

func formatMarkdown(total *pricing.TotalCost) string {
	var sb strings.Builder
	sb.WriteString("## 💰 Estimated Hetzner Cloud Costs\n\n")
	sb.WriteString("| Resource | Monthly Cost | Source |\n")
	sb.WriteString("| :--- | :--- | :--- |\n")

	for _, res := range total.Resources {
		sb.WriteString(fmt.Sprintf("| %s | %.2f EUR | %s |\n", res.Address, res.MonthlyCost, res.Source))
	}

	sb.WriteString("\n---\n")
	sb.WriteString(fmt.Sprintf("**Total Estimated Monthly Cost: %.2f EUR**\n", total.TotalMonthly))
	sb.WriteString(fmt.Sprintf("**Total Estimated Hourly Cost: %.4f EUR**\n", total.TotalHourly))
	sb.WriteString("\n*Sent by fsz-codeshop/infracost-hetzner*")

	return sb.String()
}
