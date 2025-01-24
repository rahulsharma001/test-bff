package attribute_handlers

import (
	message_constants "cee-bff-go/internal/config/constants"
	"cee-bff-go/internal/config/endpoints"
	common_model "cee-bff-go/internal/models"
	"cee-bff-go/internal/models/v1/attribute_model"
	"cee-bff-go/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetUserAttributes(c *gin.Context) {
	//Getting Table Header
	userAttributeURL := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPIGetClientAttributes

	userAttributeRequest := attribute_model.UserAttributePayload{
		PageURI:     "/data/all_contacts",
		ContactType: "anonymous",
	}

	var response attribute_model.AttributeOrderResponse
	// Create HTTP client and execute request
	err := utils.CallNetCoreAPI(userAttributeURL, userAttributeRequest, &response)
	if err != nil {
		utils.Log("Error calling Netcore API:" + err.Error())
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Error calling Netcore API.",
		})
		return
	}

	utils.SetResponse(c, utils.CommonResponse{
		StatusCode: http.StatusOK,
		Status:     message_constants.StatusSuccess,
		Message:    "Fetched Attributes Successfully.",
		Data:       response.Data.GetAttributesWithoutLongText(),
	})
}

func SaveAttributeOrder(c *gin.Context) {
	// Bind the JSON request body to the struct
	var saveAttributeRequest attribute_model.AttributeSlice
	if err := c.ShouldBind(&saveAttributeRequest); err != nil {
		utils.Log("Failed to bind request:")
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Invalid request format",
			Error:      err.Error(),
		})
		return
	}

	var saveAttributeNetcoreRequest attribute_model.UpdateAttributeOrderRequest = attribute_model.UpdateAttributeOrderRequest{
		PageURI:     "/data/all_contacts",
		PageConfigs: saveAttributeRequest,
	}

	var saveAttribtueNetcoreRequestAPI string = viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPIUpdateClientAttributes
	var saveAttributeNetcoreAPIResponse common_model.NetcoreAPIResponse
	err := utils.CallNetCoreAPI(saveAttribtueNetcoreRequestAPI, saveAttributeNetcoreRequest, &saveAttributeNetcoreAPIResponse)
	if err != nil {
		utils.Log("Error calling Netcore API:" + err.Error())
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Error calling Netcore API.",
		})
		return
	}

	utils.SetResponse(c, utils.CommonResponse{
		StatusCode: http.StatusOK,
		Status:     message_constants.StatusSuccess,
		Message:    "Saved Attributes Successfully.",
		Data:       saveAttributeNetcoreAPIResponse.Status,
	})

}
