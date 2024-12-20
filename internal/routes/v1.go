package routes

import (
	brandkit_handler "cee-bff-go/internal/handlers/v1/brandkit"
	font_handler "cee-bff-go/internal/handlers/v1/fonts"
	"cee-bff-go/internal/middleware/v1"

	"github.com/gin-gonic/gin"
)

func SetupRouterV1(version *gin.RouterGroup) {

	brandkit := version.Group("/brandkit")
	brandkit.Use(middleware.TokenValidatorMiddleware())
	brandkit.Use(middleware.DemoMiddleware())
	{
		brandkit.POST("/", brandkit_handler.GetList)
		brandkit.GET("/count", brandkit_handler.Count)
		brandkit.GET("/:id", brandkit_handler.Get)
		brandkit.GET("/active", brandkit_handler.GetActive)
		brandkit.POST("/create", brandkit_handler.Create)
		brandkit.POST("/edit", brandkit_handler.Edit)
		brandkit.DELETE("/delete/:id", brandkit_handler.Delete)
		brandkit.PATCH("/activate/:id", brandkit_handler.Activate)
		brandkit.GET("/search", brandkit_handler.Search)
		brandkit.POST("/copy/:id", brandkit_handler.Duplicate)
	}

	fonts := version.Group("/fonts")
	fonts.Use(middleware.TokenValidatorMiddleware())
	fonts.Use(middleware.DemoMiddleware())
	{
		fonts.POST("/create", font_handler.Create)
		fonts.GET("/", font_handler.GetList)
	}
}
