package pricing

import (
	"fmt"
	"github.com/fsz-codeshop/infracost-hetzner/pkg/terraform"
)

type ResourceCost struct {
	Name         string
	Address      string
	ResourceType string
	MonthlyCost  float64
	HourlyCost   float64
	Source       string
}

type TotalCost struct {
	Resources      []ResourceCost
	TotalMonthly   float64
	TotalHourly    float64
	Currency       string
	SummaryMessage string
}

func CalculateTotal(plan *terraform.Plan, engine *Engine) (*TotalCost, error) {
	total := &TotalCost{
		Currency: "EUR", // Hetzner standard
	}

	for _, change := range plan.ResourceChanges {
		// We only care about creates or updates for now
		isRelevant := false
		for _, action := range change.Change.Actions {
			if action == "create" || action == "update" || action == "no-op" {
				isRelevant = true
				break
			}
		}

		if !isRelevant {
			continue
		}

		price, err := engine.Calculate(change.Type, change.Change.After)
		if err != nil {
			// Log error but continue with other resources
			fmt.Printf("Skipping resource %s: %v\n", change.Address, err)
			continue
		}

		resCost := ResourceCost{
			Name:         change.Name,
			Address:      change.Address,
			ResourceType: change.Type,
			MonthlyCost:  price.Monthly,
			HourlyCost:   price.Hourly,
			Source:       price.Source,
		}

		total.Resources = append(total.Resources, resCost)
		total.TotalMonthly += price.Monthly
		total.TotalHourly += price.Hourly
	}

	return total, nil
}
