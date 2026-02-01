package terraform

import (
	"encoding/json"
	"os"
)

type Plan struct {
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type ResourceChange struct {
	Address      string `json:"address"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Change       Change `json:"change"`
}

type Change struct {
	Actions []string               `json:"actions"`
	Before  map[string]interface{} `json:"before"`
	After   map[string]interface{} `json:"after"`
}

func ParsePlan(path string) (*Plan, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var plan Plan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, err
	}

	return &plan, nil
}
