package contacts_handlers

import (
	message_constants "cee-bff-go/internal/config/constants"
	"cee-bff-go/internal/config/endpoints"
	common_model "cee-bff-go/internal/models"
	"cee-bff-go/internal/models/v1/attribute_model"
	"cee-bff-go/internal/models/v1/contacts_model"
	"cee-bff-go/internal/utils"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetContacts(c *gin.Context) {
	// Bind the incoming request body to the GetContactsRequest struct
	var request contacts_model.GetContactsRequest
	if err := c.BindJSON(&request); err != nil {
		utils.Log(fmt.Sprintf("Bad Request: %s", err))
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Bad Request",
		})
		return
	}

	var wg sync.WaitGroup
	var netcoreAPICountResponse, netcoreAPIDataResponse contacts_model.NetcoreContactResponse
	var err1, err2, err3 error

	errChan := make(chan error, 3)
	wg.Add(3)

	userAttributeURL := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPIGetClientAttributes

	userAttributeRequest := attribute_model.UserAttributePayload{
		PageURI:     "/data/all_contacts",
		ContactType: "anonymous",
	}

	var tableHeaderResponse attribute_model.AttributeOrderResponse

	go func() {
		defer wg.Done()
		err3 = utils.CallNetCoreAPI(userAttributeURL, userAttributeRequest, &tableHeaderResponse)
		if err3 != nil {
			errChan <- fmt.Errorf("error while calling Netcore API for data: %s", err3)
		}
	}()

	request.Fields = append(request.Fields, "CONTACT_ID")

	// Getting Count and User Data
	// API endpoint to search contacts
	searchContactsURL := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPISearchContacts

	// Get all filtering criteria from the unified function
	filteringCriteria, err := request.GetFilteringCriteria()
	if err != nil {
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Something went wrong.",
		})
		return
	}

	var filteringCriteriaOperator string
	if filteringCriteria[0].ConditionOperator == "or" {
		filteringCriteriaOperator = "and"
	} else {
		filteringCriteriaOperator = "or"
	}

	netcoreAPICountRequest := contacts_model.NetcoreContactRequest{
		FilteringCriteriaOperator: filteringCriteriaOperator,
		FilteringCriteria:         filteringCriteria,
		Output: common_model.OutputFields{
			GetCount: true,
		},
	}

	// Goroutine to make the API call for contact count
	go func() {
		defer wg.Done()
		err1 = utils.CallNetCoreAPI(searchContactsURL, netcoreAPICountRequest, &netcoreAPICountResponse)
		if err1 != nil {
			errChan <- fmt.Errorf("error while calling Netcore API for count: %s", err1)
		}
	}()

	// Prepare the request to get the contact data
	netcoreAPIDataRequest := netcoreAPICountRequest
	netcoreAPIDataRequest.Output = common_model.OutputFields{
		Sorting: []*common_model.Sorting{
			{
				Field:     request.Field,
				Direction: request.Direction,
			},
		},
		Pagination: &common_model.Pagination{
			Page:  request.Page,
			Limit: request.Limit,
		},
		Fields: request.Fields,
	}

	// Goroutine to make the API call for contact data
	go func() {
		defer wg.Done()
		err2 = utils.CallNetCoreAPI(searchContactsURL, netcoreAPIDataRequest, &netcoreAPIDataResponse)
		if err2 != nil {
			errChan <- fmt.Errorf("error while calling Netcore API for data: %s", err2)
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			// Log the error and send a bad request response
			utils.Log(fmt.Sprintf("Error: %v", err))
			utils.SetResponse(c, utils.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Status:     message_constants.StatusFailure,
				Message:    "Something Went Wrong.",
			})
			return
		}
	}

	// Extract the count and data from the API responses
	userCount := netcoreAPICountResponse.Count

	// Process the data to rename keys dynamically
	keyMapping := map[string]string{
		"CONTACT_ID": "id",
	}

	utils.RenameMultipleKeys(netcoreAPIDataResponse.Data, keyMapping)

	netcoreAPIDataResponse.ReplaceGUIDValues()

	// Create the response with the fetched count and data
	dataResponse := contacts_model.GetContactAPIResponse{
		Count:       userCount,
		Results:     netcoreAPIDataResponse.Data,
		TableHeader: tableHeaderResponse.Data.GetSelectedAttributes().GetAttributesWithTableHeaders(),
	}

	// Send the response back to the client
	utils.SetResponse(c, utils.CommonResponse{
		StatusCode: http.StatusOK,
		Status:     message_constants.StatusSuccess,
		Message:    "Data Fetched Successfully.",
		Data:       dataResponse,
	})
}

func GetContactsDropdown(c *gin.Context) {
	guid, exists := c.GetQuery("search_for")
	if !exists {
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Bad Request",
		})
		return
	}

	var netcoreAPIDataResponse contacts_model.NetcoreContactResponse

	// Prepare the request to get the contact data
	netcoreAPIDataRequest := contacts_model.NetcoreContactRequest{
		FilteringCriteriaOperator: "or",
		FilteringCriteria: []common_model.FilteringCriteria{
			{
				ConditionOperator: "and",
				ConditionDetails: []common_model.ConditionDetail{
					{
						Field:         "contact_type",
						FieldCategory: "config",
						Operation:     "equals",
						Value:         []string{"anonymous"},
					},
					{
						Field:         "GUID",
						FieldCategory: "config",
						Operation:     "contains",
						Value:         []string{guid},
					},
				},
			},
		},
		Output: common_model.OutputFields{
			GetCount: false,
			Fields: []string{
				"CONTACT_ID",
				"EMAIL",
				"GUID",
				"MOBILE",
			},
			Pagination: &common_model.Pagination{
				Page:  1,
				Limit: 10,
			},
		},
	}

	searchContactsURL := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPISearchContacts

	err := utils.CallNetCoreAPI(searchContactsURL, netcoreAPIDataRequest, &netcoreAPIDataResponse)
	if err != nil {
		// Log the error and send a bad request response
		utils.Log(fmt.Sprintf("Error: %s", err))
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Something went wrong.",
		})
		return
	}

	// Process the data to rename keys dynamically
	keyMapping := map[string]string{
		"CONTACT_ID": "id",
		"GUID":       "name",
	}

	utils.RenameMultipleKeys(netcoreAPIDataResponse.Data, keyMapping)

	responseMap := map[string]interface{}{
		"count":   len(netcoreAPIDataResponse.Data),
		"results": netcoreAPIDataResponse.Data,
	}

	// Send the response back to the client
	utils.SetResponse(c, utils.CommonResponse{
		StatusCode: http.StatusOK,
		Status:     message_constants.StatusSuccess,
		Message:    "Data Fetched Successfully.",
		Data:       responseMap,
	})
}

func GetUserProfile(c *gin.Context) {

	userID, exists := c.GetQuery("user_id")
	if !exists {
		// utils.Log(fmt.Sprintf("Bad Request: %s"))
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Bad Request",
		})
		return
	}

	var netcoreAPIDataResponse contacts_model.NetcoreContactResponse

	type UserAttributePayload struct {
		PageURI     string `json:"page_uri"`
		ContactType string `json:"contact_type"`
	}

	userAttributeRequest := UserAttributePayload{
		PageURI:     "/data/all_contacts",
		ContactType: "anonymous",
	}

	var tableHeaderResponse attribute_model.AttributeOrderResponse
	userAttributeURL := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPIGetClientAttributes
	utils.CallNetCoreAPI(userAttributeURL, userAttributeRequest, &tableHeaderResponse)

	// Prepare the request to get the contact data
	netcoreAPIDataRequest := contacts_model.NetcoreContactRequest{
		FilteringCriteriaOperator: "or",
		FilteringCriteria: []common_model.FilteringCriteria{
			{
				ConditionOperator: "and",
				ConditionDetails: []common_model.ConditionDetail{
					{
						Field:         "contact_type",
						FieldCategory: "config",
						Operation:     "equals",
						Value:         []string{"anonymous"},
					},
					{
						Field:         "contact_id",
						FieldCategory: "config",
						Operation:     "equals",
						Value:         []string{userID},
					},
				},
			},
		},
		Output: common_model.OutputFields{
			GetCount: false,
			Fields:   tableHeaderResponse.Data.GetAttributesNames(),
			Pagination: &common_model.Pagination{
				Page:  1,
				Limit: 1,
			},
		},
	}

	searchContactsURL := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPISearchContacts

	err := utils.CallNetCoreAPI(searchContactsURL, netcoreAPIDataRequest, &netcoreAPIDataResponse)
	if err != nil {
		// Log the error and send a bad request response
		utils.Log(fmt.Sprintf("Error: %s", err))
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Something went wrong.",
		})
		return
	}

	// Send the response back to the client
	utils.SetResponse(c, utils.CommonResponse{
		StatusCode: http.StatusOK,
		Status:     message_constants.StatusSuccess,
		Message:    "Data Fetched Successfully.",
		Data:       netcoreAPIDataResponse.FilterNAValues(),
	})
}

func GetUserHistory(c *gin.Context) {

	userID, exists := c.GetQuery("user_id")
	if !exists {
		// utils.Log(fmt.Sprintf("Bad Request: %s"))
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Bad Request",
		})
		return
	}

	var netcoreAPIDataResponse contacts_model.NetcoreContactResponse

	// Prepare the request to get the contact data
	netcoreAPIDataRequest := contacts_model.NetcoreContactRequest{
		FilteringCriteriaOperator: "or",
		FilteringCriteria: []common_model.FilteringCriteria{
			{
				ConditionOperator: "and",
				ConditionDetails: []common_model.ConditionDetail{
					{
						Field:         "contact_type",
						FieldCategory: "config",
						Operation:     "equals",
						Value:         []string{"anonymous"},
					},
					{
						Field:         "contact_id",
						FieldCategory: "config",
						Operation:     "equals",
						Value:         []string{userID},
					},
				},
			},
		},
		Output: common_model.OutputFields{
			Fields: []string{
				"CREATED_AT",
				"MOBILE",
			},
			GetCount: false,
			Pagination: &common_model.Pagination{
				Page:  1,
				Limit: 1,
			},
		},
	}

	searchContactsURL := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPISearchContacts

	err := utils.CallNetCoreAPI(searchContactsURL, netcoreAPIDataRequest, &netcoreAPIDataResponse)

	fmt.Println(netcoreAPIDataResponse)

	if err != nil {
		// Log the error and send a bad request response
		utils.Log(fmt.Sprintf("Error: %s", err))
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Something went wrong.",
		})
		return
	}

	// Process the data to rename keys dynamically
	keyMapping := map[string]string{
		"CREATED_AT": "date_added",
	}

	utils.RenameMultipleKeys(netcoreAPIDataResponse.Data, keyMapping)

	// Send the response back to the client
	utils.SetResponse(c, utils.CommonResponse{
		StatusCode: http.StatusOK,
		Status:     message_constants.StatusSuccess,
		Message:    "Data Fetched Successfully.",
		Data:       netcoreAPIDataResponse.Data[0],
	})
}

func GetUserSegmentList(c *gin.Context) {
	userID, exists := c.GetPostForm("user_id")
	if !exists {
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Bad Request",
		})
		return
	}

	var netcoreAPIDataResponse contacts_model.NetcoreAudienceSearchResponse

	// Prepare the request to get the contact segment
	netcoreAPIDataRequest := common_model.FilterRequest{
		AudienceType:              "segment",
		FilteringCriteriaOperator: "or",
		FilteringCriteria: []common_model.FilteringCriteria{
			{
				ConditionOperator: "and",
				ConditionDetails: []common_model.ConditionDetail{
					{
						Field:         "contact_type",
						FieldCategory: "config",
						Operation:     "equals",
						Value:         []string{"anonymous"},
					},
					{
						Field:         "contact_id",
						FieldCategory: "config",
						Operation:     "equals",
						Value:         []string{userID},
					},
				},
			},
		},
		Output: common_model.OutputFields{
			Fields: []string{
				"audience_id",
				"audience_name",
			},
			GetCount: false,
			Pagination: &common_model.Pagination{
				Page:  1,
				Limit: 1,
			},
		},
	}

	contactAudienceListAPI := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPIGetAudienceOfContacts

	err := utils.CallNetCoreAPI(contactAudienceListAPI, netcoreAPIDataRequest, &netcoreAPIDataResponse)
	if err != nil {
		// Log the error and send a bad request response
		utils.Log(fmt.Sprintf("Error: %s", err))
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     message_constants.StatusFailure,
			Message:    "Something went wrong.",
		})
		return
	}

	var response map[string]interface{} = make(map[string]interface{})
	response["segment_count"] = len(netcoreAPIDataResponse.Data)
	response["segments"] = netcoreAPIDataResponse.GetAudienceSegmentMap()

	// Send the response back to the client
	utils.SetResponse(c, utils.CommonResponse{
		StatusCode: http.StatusOK,
		Status:     message_constants.StatusSuccess,
		Message:    "Data Fetched Successfully.",
		Data:       response,
	})
}
