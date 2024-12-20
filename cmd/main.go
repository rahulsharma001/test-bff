package main

import (
	"cee-bff-go/internal/routes"

	"cee-bff-go/internal/utils"

	"github.com/gin-gonic/gin"
)

// SetupServer initializes the Gin engine and sets up routes
func SetupServer() *gin.Engine {
	server := gin.Default()
	routes.SetupRouter(server)
	return server
}

func main() {

	//Initializing config
	utils.InitializeConfig()

	server := SetupServer()
	server.Run(":8080")

}
