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
	Count      *int        `json:"count,omitempty"`
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
		Count:      commonResp.Count,
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
