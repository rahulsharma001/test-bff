package workers

import (
	"fmt"

	"github.com/conductor-sdk/conductor-go/sdk/model"
)

// MapPayloadTask is the function that implements the custom task logic for payload mapping
func MapPayloadTask(task *model.Task) (interface{}, error) {
	// Extract the raw payload from task input
	rawPayload, ok := task.InputData["raw_payload"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid raw payload format")
	}

	// Perform payload mapping
	mappedPayload := map[string]interface{}{
		"data": map[string]interface{}{
			"name":    rawPayload["name"],
			"details": rawPayload["details"],
		},
	}

	// Return the mapped payload
	return map[string]interface{}{
		"mapped_payload": mappedPayload,
	}, nil
}
