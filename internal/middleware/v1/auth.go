package middleware

import (
	"cee-bff-go/internal/config/endpoints"
	auth_model "cee-bff-go/internal/models/v1/auth"
	"cee-bff-go/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type AuthAPIResponse struct {
	RequestID string `json:"request_id"`
	Status    int64  `json:"status"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Errors    string `json:"errors,omitempty"`
	Data      string `json:"data,omitempty"`
}

func TokenValidatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Skip token validation for preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			utils.Log("Authorization Not Found in Header.")
			utils.SetResponse(c, utils.CommonResponse{
				StatusCode: http.StatusUnauthorized,
				Status:     "failure",
				Message:    "Token mismatched",
				Data:       nil,
			})
			c.Abort()
			return
		}

		type Payload struct {
			Identifier string `json:"identifier"`
			Secret     string `json:"secret"`
		}

		var payload = Payload{
			Identifier: viper.GetString("CLIENT_IDENTIFIER"),
			Secret:     viper.GetString("CLIENT_SECRET"),
		}

		// Marshal the payload back into a JSON string for the HTTP request body
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			utils.Log("Error While Marshalling and creating request")
			utils.SetResponse(c, utils.CommonResponse{
				StatusCode: http.StatusUnauthorized,
				Status:     "failure",
				Message:    "Token mismatched",
				Data:       nil,
			})
			c.Abort()
			return
		}

		// Use the payloadJSON as your payload in the HTTP request

		headers := map[string]string{
			"Accept":        "application/json",
			"Content-Type":  "application/json",
			"Authorization": authorization,
		}

		url := viper.GetString("NCS_ENDPOINT") + endpoints.NCSValidateToken

		httpParams := utils.NewHTTP("POST", url, headers, string(payloadJSON))

		response, err := httpParams.DoHTTP()
		fmt.Println(err)

		if err != nil {
			utils.Log("Something went wrong while validating authorization.")
			utils.SetResponse(c, utils.CommonResponse{
				StatusCode: http.StatusUnauthorized,
				Status:     "error",
				Message:    "Something went wrong while validating authorization.",
				Data:       nil,
			})
			c.Abort()
			return
		}

		var authAPIResponse auth_model.Response
		err = json.Unmarshal([]byte(response), &authAPIResponse)
		if err != nil {
			utils.Log("Something went wrong while validating authorization.")
			utils.SetResponse(c, utils.CommonResponse{
				StatusCode: http.StatusUnauthorized,
				Status:     "error",
				Message:    "Something went wrong while validating authorization.",
				Data:       nil,
			})
			c.Abort()
			return
		}

		if authAPIResponse.Data.APIKey == "" {
			utils.Log("Empty API key recived from NCS for this user.")
			utils.SetResponse(c, utils.CommonResponse{
				StatusCode: http.StatusUnauthorized,
				Status:     "error",
				Message:    "Something went wrong while validating authorization.",
				Data:       nil,
			})
			c.Abort()
			return
		}

		viper.Set("NETCORE_API_KEY", authAPIResponse.Data.APIKey)

		c.Next()
	}
}
