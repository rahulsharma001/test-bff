package routes

import (
	attribute_handlers "cee-bff-go/internal/handlers/v1/attribute_handlers"
	contacts_handlers "cee-bff-go/internal/handlers/v1/contacts_handlers"
	"cee-bff-go/internal/middleware/v1"

	"github.com/gin-gonic/gin"
)

func SetupRouterV1(version *gin.RouterGroup) {

	// Attribute routes
	attributeRoutes := version.Group("/attributes")
	{
		attributeRoutes.Use(middleware.DemoMiddleware())
		attributeRoutes.Use(middleware.EndpointTimer())
		attributeRoutes.Use(middleware.TokenValidatorMiddleware())
		attributeRoutes.GET("/get_attributes", attribute_handlers.GetUserAttributes)
		attributeRoutes.POST("/save_attributes_order", attribute_handlers.SaveAttributeOrder)
	}

	// Contact routes
	contactRoutes := version.Group("/contacts")
	{
		contactRoutes.Use(middleware.DemoMiddleware())
		contactRoutes.Use(middleware.EndpointTimer())
		contactRoutes.Use(middleware.TokenValidatorMiddleware())
		contactRoutes.POST("/get_contacts", contacts_handlers.GetContacts)
		contactRoutes.GET("/get_contacts", contacts_handlers.GetContactsDropdown)
		contactRoutes.GET("/get_user_history", contacts_handlers.GetUserHistory)
		contactRoutes.GET("/get_user_profile", contacts_handlers.GetUserProfile)
		contactRoutes.POST("/get_user_list_segments", contacts_handlers.GetUserSegmentList)
	}

}
