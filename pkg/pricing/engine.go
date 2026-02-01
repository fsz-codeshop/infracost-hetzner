package pricing

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed data/fallback_prices.json
var fallbackPricesData []byte

type FallbackRecord struct {
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	MonthlyNet float64 `json:"monthly_net"`
	HourlyNet  float64 `json:"hourly_net"`
}

type PriceInfo struct {
	Monthly float64
	Hourly  float64
	Source  string
	Date    time.Time
}

type PriceProvider interface {
	GetPrice(resourceType string, attributes map[string]interface{}) (*PriceInfo, error)
}

type Engine struct {
	API      *HcloudAPIProvider
	Fallback *FallbackProvider
}

func (e *Engine) Calculate(resourceType string, attributes map[string]interface{}) (*PriceInfo, error) {
	// 1. Try API first
	if e.API != nil {
		price, err := e.API.GetPrice(resourceType, attributes)
		if err == nil {
			return price, nil
		}
		fmt.Printf("Warning: Could not fetch price from API for %s, trying fallback: %v\n", resourceType, err)
	}

	// 2. Fallback
	if e.Fallback != nil {
		return e.Fallback.GetPrice(resourceType, attributes)
	}

	return nil, fmt.Errorf("no price provider available for resource type: %s", resourceType)
}

// HcloudAPIProvider fetches prices via official SDK
type HcloudAPIProvider struct {
	Client *hcloud.Client
}

func (h *HcloudAPIProvider) GetPrice(resourceType string, attributes map[string]interface{}) (*PriceInfo, error) {
	ctx := context.Background()

	switch resourceType {
	case "hcloud_server":
		stName, _ := attributes["server_type"].(string)
		location, _ := attributes["location"].(string)
		if stName == "" {
			return nil, fmt.Errorf("missing server_type")
		}

		st, _, err := h.Client.ServerType.GetByName(ctx, stName)
		if err != nil {
			return nil, err
		}
		if st == nil {
			return nil, fmt.Errorf("server type %s not found", stName)
		}

		// Find pricing for the location
		for _, p := range st.Pricings {
			if location == "" || p.Location.Name == location {
				monthly, _ := strconv.ParseFloat(p.Monthly.Net, 64)
				hourly, _ := strconv.ParseFloat(p.Hourly.Net, 64)
				return &PriceInfo{
					Monthly: monthly,
					Hourly:  hourly,
					Source:  "Hetzner API",
					Date:    time.Now(),
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("unsupported resource type or location in API provider")
}

// FallbackProvider handles prices from embedded JSON
type FallbackProvider struct {
	FallbackDate time.Time
	Prices       map[string]map[string]FallbackRecord // Type -> Name -> Record
}

// NewFallbackProvider initializes the provider with embedded data
func NewFallbackProvider() (*FallbackProvider, error) {
	var records []FallbackRecord
	if err := json.Unmarshal(fallbackPricesData, &records); err != nil {
		return nil, fmt.Errorf("failed to parse fallback prices: %v", err)
	}

	prices := make(map[string]map[string]FallbackRecord)
	for _, rec := range records {
		if _, ok := prices[rec.Type]; !ok {
			prices[rec.Type] = make(map[string]FallbackRecord)
		}
		prices[rec.Type][rec.Name] = rec
	}

	return &FallbackProvider{
		FallbackDate: time.Now(), // Ideally this would come from a version/date file
		Prices:       prices,
	}, nil
}

func (f *FallbackProvider) GetPrice(resourceType string, attributes map[string]interface{}) (*PriceInfo, error) {
	switch resourceType {
	case "hcloud_server":
		serverType, _ := attributes["server_type"].(string)
		if serverType == "" {
			return nil, fmt.Errorf("missing server_type")
		}

		if typeGroup, ok := f.Prices["server"]; ok {
			if rec, ok := typeGroup[serverType]; ok {
				return &PriceInfo{
					Monthly: rec.MonthlyNet,
					Hourly:  rec.HourlyNet,
					Source:  "Fallback (Embedded)",
					Date:    f.FallbackDate,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("resource not in fallback table")
}
