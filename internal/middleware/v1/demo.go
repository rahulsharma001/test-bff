package middleware

import (
	"cee-bff-go/internal/demo"
	"cee-bff-go/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware to handle demo mode using the demo package's DemoResponses map
func DemoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		isDemoMode := c.Request.Header.Get("Demo-Panel")

		if isDemoMode == "1" {
			demo.InitializeDemoData(c)
			// Check if there is a dummy response for the current endpoint
			if dummyData, exists := demo.DemoResponses[c.FullPath()]; exists {
				var result map[string]interface{}

				// Parse the JSON data into the map
				err := json.Unmarshal([]byte(dummyData), &result)
				if err != nil {
					utils.Log(err)
				}

				// Return the dummy data and do not proceed to the actual handler
				c.JSON(http.StatusOK, result)
				c.Abort() // Prevent the next handlers from running
				return
			}
		}

		// If not in demo mode or no dummy data found, continue to the actual handler
		c.Next()
	}
}
