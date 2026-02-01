package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fsz-codeshop/infracost-hetzner/pkg/config"
	"github.com/fsz-codeshop/infracost-hetzner/pkg/github"
	"github.com/fsz-codeshop/infracost-hetzner/pkg/pricing"
	"github.com/fsz-codeshop/infracost-hetzner/pkg/terraform"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "infracost-hetzner",
	Short: "Estimate Hetzner Cloud costs from Terraform plans",
	Long:  `A CLI tool to parse Terraform plan JSON and estimate costs using Hetzner Cloud API.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		slog.SetDefault(logger)

		// 1. Load Config
		cfg, err := config.LoadConfig(cmd)
		if err != nil {
			logger.Error("Failed to load configuration", "error", err)
			os.Exit(1)
		}

		logger.Info("Starting cost estimate", "plan_path", cfg.PlanPath)

		// 2. Parse Plan
		plan, err := terraform.ParsePlan(cfg.PlanPath)
		if err != nil {
			logger.Error("Failed to parse plan", "path", cfg.PlanPath, "error", err)
			os.Exit(1)
		}
		logger.Info("Plan parsed successfully", "resource_changes", len(plan.ResourceChanges))

		// 3. Setup Pricing Engine
		fallback, err := pricing.NewFallbackProvider()
		if err != nil {
			logger.Error("Failed to load fallback prices", "error", err)
			os.Exit(1)
		}

		engine := &pricing.Engine{
			Fallback: fallback,
		}

		if cfg.HcloudToken != "" {
			logger.Info("Initializing Hcloud API client")
			engine.API = &pricing.HcloudAPIProvider{
				Client: hcloud.NewClient(hcloud.WithToken(cfg.HcloudToken)),
			}
		} else {
			logger.Warn("HCLOUD_TOKEN not set, using fallback pricing only")
		}

		// 4. Calculate
		total, err := pricing.CalculateTotal(plan, engine)
		if err != nil {
			logger.Error("Failed to calculate costs", "error", err)
			os.Exit(1)
		}

		// 5. Output to Console (Keep this as direct formatted output for human readability in CI logs)
		fmt.Printf("\n💰 Estimated Hetzner Cloud Costs\n")
		fmt.Printf("---------------------------------\n")
		for _, res := range total.Resources {
			fmt.Printf("%-40s | %8.2f EUR/mo | (%s)\n", res.Address, res.MonthlyCost, res.Source)
		}
		fmt.Printf("---------------------------------\n")
		fmt.Printf("Total Monthly: %.2f EUR\n", total.TotalMonthly)
		fmt.Printf("Total Hourly:  %.4f EUR\n\n", total.TotalHourly)

		// 6. Output to GitHub
		if cfg.GithubToken != "" {
			logger.Info("Posting comment to GitHub PR")
			err = github.CommentPR(total, cfg)
			if err != nil {
				logger.Error("Failed to post GitHub comment", "error", err)
			} else {
				logger.Info("GitHub comment posted successfully")
			}
		} else {
			logger.Info("Skipping GitHub comment (GITHUB_TOKEN not set)")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("plan", "p", "", "Path to the Terraform plan JSON file")
	rootCmd.PersistentFlags().StringP("token", "t", "", "Hetzner Cloud API Token")
}
