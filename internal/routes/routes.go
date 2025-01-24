package routes

import (
	message_constants "cee-bff-go/internal/config/constants"
	"cee-bff-go/internal/middleware/v1"
	"cee-bff-go/internal/utils"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetupRouter(server *gin.Engine) {

	// Apply CORS middleware
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{viper.GetString("ORIGIN_HOST")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Requested-With", "sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-platform", "Demo-Panel"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}))
	server.GET("/healthcheck", healthcheck)

	server.Use(middleware.RequestIDMiddleware())

	version := server.Group("/v1")
	{
		SetupRouterV1(version)
	}

}

func healthcheck(c *gin.Context) {
	utils.SetResponse(c, utils.CommonResponse{
		StatusCode: http.StatusOK,
		Status:     message_constants.StatusSuccess,
	})
	return
}
