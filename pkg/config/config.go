package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	PlanPath     string
	HcloudToken  string
	GithubToken  string
	GithubRepo   string
	PRNumber     string
	Debug        bool
}

func LoadConfig(cmd *cobra.Command) (*Config, error) {
	planPath, _ := cmd.Flags().GetString("plan")
	if planPath == "" {
		return nil, fmt.Errorf("plan path is required")
	}

	cfg := &Config{
		PlanPath:    planPath,
		HcloudToken: os.Getenv("HCLOUD_TOKEN"),
		GithubToken: os.Getenv("GITHUB_TOKEN"),
		GithubRepo:  os.Getenv("GITHUB_REPOSITORY"),
		PRNumber:    os.Getenv("PR_NUMBER"),
	}

	// Fallback for token from flags if provided
	if cfg.HcloudToken == "" {
		cfg.HcloudToken, _ = cmd.Flags().GetString("token")
	}

	return cfg, nil
}
