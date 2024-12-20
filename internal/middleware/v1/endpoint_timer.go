package middleware

import (
	"cee-bff-go/internal/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var slowAPIThreshold = 1

// Middleware to handle demo mode using the demo package's DemoResponses map
func EndpointTimer() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// Proceed to the next handler
		c.Next()
		elapsed := time.Since(start)
		if elapsed > (time.Duration(slowAPIThreshold) * time.Second) {
			utils.SendSlackAlert(viper.GetString("SLACK_WEBHOOK_URL"), fmt.Sprintf("Slow API Alert : API request for %s took %s", c.Request.URL.Path, elapsed))
		}
	}
}
