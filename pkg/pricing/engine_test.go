package pricing

import (
	"testing"
)

func TestFallbackProvider_GetPrice(t *testing.T) {
	// Initialize with embedded data
	provider, err := NewFallbackProvider()
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	tests := []struct {
		name         string
		resourceType string
		attributes   map[string]interface{}
		expectError  bool
		expectedCost float64
	}{
		{
			name:         "Valid server cx11",
			resourceType: "hcloud_server",
			attributes:   map[string]interface{}{"server_type": "cx11"},
			expectError:  false,
			expectedCost: 3.79,
		},
		{
			name:         "Valid server ccx63",
			resourceType: "hcloud_server",
			attributes:   map[string]interface{}{"server_type": "ccx63"},
			expectError:  false,
			expectedCost: 919.00,
		},
		{
			name:         "Invalid resource type",
			resourceType: "invalid_type",
			attributes:   map[string]interface{}{},
			expectError:  true,
		},
		{
			name:         "Unknown server type",
			resourceType: "hcloud_server",
			attributes:   map[string]interface{}{"server_type": "unknown-cpu"},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := provider.GetPrice(tt.resourceType, tt.attributes)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !tt.expectError && price != nil {
				if price.Monthly != tt.expectedCost {
					t.Errorf("expected monthly cost: %f, got: %f", tt.expectedCost, price.Monthly)
				}
				if price.Source != "Fallback (Embedded)" {
					t.Errorf("expected source Fallback (Embedded), got %s", price.Source)
				}
			}
		})
	}
}

func TestCalculate(t *testing.T) {
	provider, _ := NewFallbackProvider()
	engine := &Engine{Fallback: provider}

	// Mock data
	reqType := "hcloud_server"
	attrs := map[string]interface{}{"server_type": "cpx11"}

	price, err := engine.Calculate(reqType, attrs)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if price.Monthly != 4.35 {
		t.Errorf("Expected 4.35, got %f", price.Monthly)
	}
}
