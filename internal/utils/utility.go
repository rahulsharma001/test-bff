package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// CommonResponse is the structure for a standard response
type CommonResponse struct {
	RequestID  string      `json:"request_id"`
	StatusCode int         `json:"code"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}

// Respond is a common function to return a JSON response
func SetResponse(c *gin.Context, commonResp CommonResponse) {
	c.JSON(commonResp.StatusCode, CommonResponse{
		RequestID:  viper.GetString("request_id"),
		StatusCode: commonResp.StatusCode,
		Status:     commonResp.Status,
		Message:    commonResp.Message,
		Data:       commonResp.Data,
		Error:      commonResp.Error,
	})
}

type SlackMessage struct {
	Text string `json:"text"`
}

func SendSlackAlert(webhookURL string, message string) error {
	payload := SlackMessage{
		Text: message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error sending Slack alert: status code %d", resp.StatusCode)
	}

	return nil
}

func CallNetCoreAPI(url string, request interface{}, destination interface{}) error {

	// Marshal request to JSON
	netcoreAPIRequestJSON, err := json.Marshal(request)
	Log(fmt.Sprintf("Calling: %s ,Created the request :  %s", url, netcoreAPIRequestJSON))
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to marshal Request :  %s", err))
	}

	// Prepare HTTP client and make the request

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
		"api-key":      viper.GetString("NETCORE_API_KEY"),
	}

	httpClient := NewHTTP("POST", url, headers, string(netcoreAPIRequestJSON))
	response, err := httpClient.DoHTTP()
	if err != nil {
		Log(fmt.Sprintf("Error received while calling netcore API:   %s. Response : %s. Request : %s. Headers : %s ", err, response, request, headers))
		return errors.New(fmt.Sprintf("Error received while calling netcore API:   %s", err))
	}

	// Unmarshal API response into the struct
	if err := json.Unmarshal([]byte(response), destination); err != nil {
		Log(fmt.Sprintf("Failed to parse netcore API response :   %s. Error : %s", response, err))
		return errors.New(fmt.Sprintf("Failed to parse netcore API response:   %s", err))
	}

	return nil
}

// RenameAndProcessData renames keys in a slice of maps based on the provided keyMapping.
// It takes two parameters: the input data (of type interface{}) and the keyMapping (a map of old keys to new keys).
func RenameAndProcessData(data interface{}, keyMapping map[string]string) interface{} {
	// Ensure that the data is a slice of interfaces (which is a []map[string]interface{})
	if slice, ok := data.([]interface{}); ok {
		for i, item := range slice {
			// Ensure that each item in the slice is a map
			if contact, ok := item.(map[string]interface{}); ok {
				// Log the contact before renaming (for debugging)
				fmt.Printf("Before renaming: %+v\n", contact)

				// Rename keys based on the keyMapping
				for oldKey, newKey := range keyMapping {
					if value, exists := contact[oldKey]; exists {
						contact[newKey] = value
						delete(contact, oldKey)
						// Log the renamed key for debugging
						fmt.Printf("Renamed key: %s -> %s\n", oldKey, newKey)
					}
				}
				// Update the item in the slice with the modified map
				slice[i] = contact
			}
		}
		// Return the modified slice
		return slice
	}
	// If the data is not a slice, return it unchanged or handle error
	return data
}

// Rename multiple keys in a slice of maps with values of type string or int
func RenameMultipleKeys(maps []map[string]interface{}, keyReplacements map[string]string) {
	for _, m := range maps {
		for oldKey, newKey := range keyReplacements {
			if value, exists := m[oldKey]; exists {
				// Set the new key with the value of the old key
				m[newKey] = value
				// Remove the old key
				delete(m, oldKey)
			}
		}
	}
}
