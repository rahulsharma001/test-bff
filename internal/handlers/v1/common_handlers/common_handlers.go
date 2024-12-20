package common_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCommonHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": "Dummy Data"})
}
