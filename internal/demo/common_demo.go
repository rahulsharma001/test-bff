package demo

import (
	"github.com/gin-gonic/gin"
)

// Exported map to hold all demo responses
var DemoResponses = make(map[string]string)

// InitializeDemoData calls all registration functions to populate DemoResponses
func InitializeDemoData(c *gin.Context) {
	RegisterBrandKitDemoData(c)
	RegisterFontsDemoData(c)
}
