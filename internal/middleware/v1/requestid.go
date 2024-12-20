package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// Middleware to generate a unique request ID and attach it to the request context
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a new UUID for the request
		requestID := uuid.New().String()

		// Set the RequestID in the request context
		viper.Set("REQUEST_ID", requestID)

		// Proceed with the request
		c.Next()
	}
}
