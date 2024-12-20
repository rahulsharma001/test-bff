package service

import (
	"cee-bff-go/internal/config/endpoints"
	brandkit_model "cee-bff-go/internal/models/v1/brandkit"
	netcoreapi_model "cee-bff-go/internal/models/v1/netcoreapi"
	"cee-bff-go/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type BrandKit brandkit_model.BrandKit

func prepareHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		//"api-key":      viper.GetString("NETCORE_API_KEY"),
		"api-key": "a472142b46f3cca779251edd21b9bf50",
	}
}

func createRequestBody(operator string, conditionOperator string, criteria []netcoreapi_model.ConditionDetail, getCount bool, page, limit int, output bool) netcoreapi_model.APIRequest {

	request := netcoreapi_model.APIRequest{
		FilteringCriteriaOperator: operator,
		FilteringCriteria: []netcoreapi_model.FilteringCriteria{
			{
				ConditionOperator: conditionOperator,
				ConditionDetails:  criteria,
			},
		},
	}

	if output {
		request.Output = &netcoreapi_model.OutputConfig{
			GetCount: getCount,
			Sorting: []netcoreapi_model.SortingOption{
				{
					Field:     "status",
					Direction: "desc",
				},
				{
					Field:     "updated_at",
					Direction: "desc",
				},
			},
			Pagination: &netcoreapi_model.Pagination{
				Page:  page,
				Limit: limit,
			},
			Fields: []string{"id", "name", "details", "status", "created_at", "updated_at"},
		}
	}
	return request
}

func sendRequest(endpoint string, jsonData string, headers map[string]string) (string, error) {
	httpParams := utils.NewHTTP("POST", endpoint, headers, jsonData)
	httpParams.Timeout = 10 * time.Second
	return httpParams.DoHTTP()
}

func GetAll(request *brandkit_model.GetList) (brandkit_model.BrandKitData, error) {
	headers := prepareHeaders()
	conditionDetails := []netcoreapi_model.ConditionDetail{}
	if request.SearchFor != "" {
		conditionDetails = append(conditionDetails, netcoreapi_model.ConditionDetail{
			Field:         "name",
			FieldCategory: "config",
			Operation:     "like",
			Value:         []interface{}{request.SearchFor},
		}, netcoreapi_model.ConditionDetail{
			Field:         "details",
			FieldCategory: "config",
			Operation:     "like",
			Value:         []interface{}{request.SearchFor},
		})
	}

	netcoreApiRequestBody := createRequestBody("or", "or", conditionDetails, false, int(request.Page), int(request.Limit), true)
	jsonData, err := json.Marshal(netcoreApiRequestBody)
	if err != nil {
		return brandkit_model.BrandKitData{}, fmt.Errorf("error marshalling JSON: %w", err)
	}

	fmt.Println(viper.GetString("NETCORE_API_ENDPOINT"))
	url := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPISearchBrandkit
	utils.Log(fmt.Sprintf("GeAll Brandkit Request: Url : %s  , request : %s  ,headers : %+v ", url, string(jsonData), headers))
	response, err := sendRequest(url, string(jsonData), headers)

	if err != nil {
		return brandkit_model.BrandKitData{}, err
	}
	//var brandKitAPIResponse brandkit_model.BrandKitAPIResponse
	var brandKitAPIResponse brandkit_model.NetcoreAPIBrandKitResponse

	//var brandKitAPIResponse2 brandkit_model.NetcoreAPIBrandKitResponse
	if err := json.Unmarshal([]byte(response), &brandKitAPIResponse); err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	var target []brandkit_model.BrandKit
	for _, kit := range brandKitAPIResponse.BrandKits {
		var tempkit brandkit_model.BrandKit
		var details brandkit_model.Details
		if err := json.Unmarshal([]byte(kit.Details), &details); err != nil {
			log.Fatal(err)
			return brandkit_model.BrandKitData{}, err
		}
		tempkit.BrandKitId = kit.BrandKitId
		tempkit.Name = kit.Name
		tempkit.CreatedAt = kit.CreatedAt
		tempkit.UpdatedAt = kit.UpdatedAt
		tempkit.Status = kit.Status
		tempkit.Details = details
		target = append(target, tempkit)
	}

	count, err := get_count(netcoreApiRequestBody, url, headers)

	if err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	return brandkit_model.BrandKitData{
		Data:  target,
		Count: &count.Count,
	}, nil
}

func get_count(netcoreApiRequestBody netcoreapi_model.APIRequest, url string, headers map[string]string) (brandkit_model.BrandKitAPIResponse, error) {
	netcoreApiRequestBody.Output.GetCount = true
	netcoreApiRequestBody.Output.Pagination = nil
	countJsonData, _ := json.Marshal(netcoreApiRequestBody)
	utils.Log(fmt.Sprintf("Brandkit count Request: Url : %s  , request : %s  ,headers : %+v ", url, string(countJsonData), headers))
	countResponse, err := sendRequest(url, string(countJsonData), headers)
	var brandKitAPIResponse brandkit_model.BrandKitAPIResponse
	if err != nil {
		return brandKitAPIResponse, err
	}

	if err := json.Unmarshal([]byte(countResponse), &brandKitAPIResponse); err != nil {
		return brandKitAPIResponse, err
	}

	return brandKitAPIResponse, nil
}

func Get(id int64) (brandkit_model.BrandKitData, error) {
	headers := prepareHeaders()
	conditionDetails := []netcoreapi_model.ConditionDetail{
		{
			Field:         "id",
			FieldCategory: "config",
			Operation:     "equals",
			Value:         []interface{}{id},
		},
	}
	netcoreApiRequestBody := createRequestBody("or", "or", conditionDetails, false, 1, 10, true)
	jsonData, err := json.Marshal(netcoreApiRequestBody)
	if err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	url := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPISearchBrandkit
	utils.Log(fmt.Sprintf("Get Brandkit Request: Url : %s  , request : %s  ,headers : %+v ", url, string(jsonData), headers))
	response, err := sendRequest(url, string(jsonData), headers)
	if err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	//var brandKitAPIResponse brandkit_model.BrandKitAPIResponse
	/* if err := json.Unmarshal([]byte(response), &brandKitAPIResponse); err != nil {
		return brandkit_model.BrandKitData{}, err
	} */

	var brandKitAPIResponse brandkit_model.NetcoreAPIBrandKitResponse
	if err := json.Unmarshal([]byte(response), &brandKitAPIResponse); err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	var target []brandkit_model.BrandKit
	for _, kit := range brandKitAPIResponse.BrandKits {
		var tempkit brandkit_model.BrandKit
		var details brandkit_model.Details
		if err := json.Unmarshal([]byte(kit.Details), &details); err != nil {
			log.Fatal(err)
		}
		tempkit.BrandKitId = kit.BrandKitId
		tempkit.Name = kit.Name
		tempkit.CreatedAt = kit.CreatedAt
		tempkit.UpdatedAt = kit.UpdatedAt
		tempkit.Status = kit.Status
		tempkit.Details = details
		target = append(target, tempkit)
	}

	count, err := get_count(netcoreApiRequestBody, url, headers)

	if err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	return brandkit_model.BrandKitData{
		Data:  target,
		Count: &count.Count,
	}, nil
}

func (brandKit *BrandKit) Create() error {
	headers := prepareHeaders()
	netcoreApiRequestBody := netcoreapi_model.APIRequest{Data: brandKit}
	byteData, err := json.Marshal(netcoreApiRequestBody)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}
	url := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPICreateBrandkit
	utils.Log(fmt.Sprintf("Create Brandkit Request: Url : %s  , request : %s  ,headers : %+v ", url, string(byteData), headers))
	_, err = sendRequest(url, string(byteData), headers)
	return err
}

func (brandKit *BrandKit) Edit() error {
	headers := prepareHeaders()
	netcoreApiRequestBody := netcoreapi_model.APIRequest{Data: brandKit}
	byteData, err := json.Marshal(netcoreApiRequestBody)
	if err != nil {
		return err
	}

	url := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPIEditBrandkit
	utils.Log(fmt.Sprintf("Edit Brandkit Request: Url : %s  , request : %s  ,headers : %+v ", url, string(byteData), headers))
	_, err = sendRequest(url, string(byteData), headers)
	return err
}

func Delete(id int64) error {
	headers := prepareHeaders()
	conditionDetails := []netcoreapi_model.ConditionDetail{
		{
			Field:         "id",
			FieldCategory: "config",
			Operation:     "equals",
			Value:         []interface{}{id},
		},
	}
	netcoreApiRequestBody := createRequestBody("or", "or", conditionDetails, false, 1, 1, false)
	jsonData, err := json.Marshal(netcoreApiRequestBody)

	if err != nil {
		return err
	}

	url := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPIDeleteBrandkit
	utils.Log(fmt.Sprintf("Delete Brandkit Request: Url : %s  , request : %s  ,headers : %+v ", url, string(jsonData), headers))
	_, err = sendRequest(url, string(jsonData), headers)
	return err
}

func Activate(id int64) error {
	headers := prepareHeaders()
	conditionDetails := []netcoreapi_model.ConditionDetail{
		{
			Field:         "id",
			FieldCategory: "config",
			Operation:     "equals",
			Value:         []interface{}{id},
		},
	}
	netcoreApiRequestBody := createRequestBody("or", "or", conditionDetails, false, 1, 1, false)
	jsonData, err := json.Marshal(netcoreApiRequestBody)
	if err != nil {
		return err
	}

	url := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPIActivateBrandkit
	utils.Log(fmt.Sprintf("Activate Brandkit Request: Url : %s  , request : %s  ,headers : %+v ", url, string(jsonData), headers))
	_, err = sendRequest(url, string(jsonData), headers)
	return err
}

func Duplicate(id int64) error {
	headers := prepareHeaders()
	conditionDetails := []netcoreapi_model.ConditionDetail{
		{
			Field:         "id",
			FieldCategory: "config",
			Operation:     "equals",
			Value:         []interface{}{id},
		},
	}
	netcoreApiRequestBody := createRequestBody("or", "or", conditionDetails, false, 1, 1, false)
	jsonData, err := json.Marshal(netcoreApiRequestBody)
	if err != nil {
		return err
	}

	url := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPICopyBrandkit
	utils.Log(fmt.Sprintf("Duplicate Brandkit Request: Url : %s  , request : %s  ,headers : %+v ", url, string(jsonData), headers))
	_, err = sendRequest(url, string(jsonData), headers)
	return err
}

func GetActive() (brandkit_model.BrandKitData, error) {
	headers := prepareHeaders()
	conditionDetails := []netcoreapi_model.ConditionDetail{
		{
			Field:         "status",
			FieldCategory: "config",
			Operation:     "equals",
			Value:         []interface{}{true},
		},
	}
	netcoreApiRequestBody := createRequestBody("or", "or", conditionDetails, false, 1, 1, true)
	fmt.Println(netcoreApiRequestBody)
	jsonData, err := json.Marshal(netcoreApiRequestBody)
	if err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	url := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPISearchBrandkit
	utils.Log(fmt.Sprintf("Search Brandkit Request: Url : %s  , request : %s  ,headers : %+v ", url, string(jsonData), headers))
	response, err := sendRequest(url, string(jsonData), headers)
	if err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	//var brandKitAPIResponse brandkit_model.BrandKitAPIResponse
	/* if err := json.Unmarshal([]byte(response), &brandKitAPIResponse); err != nil {
		return brandkit_model.BrandKitData{}, err
	} */

	var brandKitAPIResponse brandkit_model.NetcoreAPIBrandKitResponse
	if err := json.Unmarshal([]byte(response), &brandKitAPIResponse); err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	var target []brandkit_model.BrandKit
	for _, kit := range brandKitAPIResponse.BrandKits {
		var tempkit brandkit_model.BrandKit
		var details brandkit_model.Details
		if err := json.Unmarshal([]byte(kit.Details), &details); err != nil {
			log.Fatal(err)
		}
		tempkit.BrandKitId = kit.BrandKitId
		tempkit.Name = kit.Name
		tempkit.CreatedAt = kit.CreatedAt
		tempkit.UpdatedAt = kit.UpdatedAt
		tempkit.Status = kit.Status
		tempkit.Details = details
		target = append(target, tempkit)
	}

	count, err := get_count(netcoreApiRequestBody, url, headers)

	if err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	return brandkit_model.BrandKitData{
		Data:  target,
		Count: &count.Count,
	}, nil
}

func Search(query string, page int, limit int) (brandkit_model.BrandKitData, error) {
	headers := prepareHeaders()
	conditionDetails := []netcoreapi_model.ConditionDetail{
		{
			Field:         "name",
			FieldCategory: "config",
			Operation:     "like",
			Value:         []interface{}{query},
		},
		{
			Field:         "details",
			FieldCategory: "config",
			Operation:     "like",
			Value:         []interface{}{query},
		},
	}
	netcoreApiRequestBody := createRequestBody("or", "or", conditionDetails, false, page, limit, true)
	jsonData, err := json.Marshal(netcoreApiRequestBody)
	if err != nil {
		return brandkit_model.BrandKitData{}, fmt.Errorf("error marshalling JSON: %w", err)
	}

	url := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPISearchBrandkit
	utils.Log(fmt.Sprintf("Search Brandkit Request: Url : %s  , request : %s  ,headers : %+v ", url, string(jsonData), headers))
	response, err := sendRequest(url, string(jsonData), headers)
	if err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	//var brandKitAPIResponse brandkit_model.BrandKitAPIResponse

	var brandKitAPIResponse2 brandkit_model.NetcoreAPIBrandKitResponse
	if err := json.Unmarshal([]byte(response), &brandKitAPIResponse2); err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	var target []brandkit_model.BrandKit
	for _, kit := range brandKitAPIResponse2.BrandKits {
		var tempkit brandkit_model.BrandKit
		var details brandkit_model.Details
		if err := json.Unmarshal([]byte(kit.Details), &details); err != nil {
			log.Fatal(err)
		}
		tempkit.BrandKitId = kit.BrandKitId
		tempkit.Name = kit.Name
		tempkit.CreatedAt = kit.CreatedAt
		tempkit.UpdatedAt = kit.UpdatedAt
		tempkit.Status = kit.Status
		tempkit.Details = details
		target = append(target, tempkit)
	}

	count, err := get_count(netcoreApiRequestBody, url, headers)

	if err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	return brandkit_model.BrandKitData{
		Data:  target,
		Count: &count.Count,
	}, nil
}

func Count() (brandkit_model.BrandKitData, error) {
	headers := prepareHeaders()

	var brandKitData brandkit_model.BrandKitData

	netcoreApiRequestBody := createRequestBody("or", "or", nil, true, 1, 1, true)
	jsonData, err := json.Marshal(netcoreApiRequestBody)
	if err != nil {
		return brandKitData, err
	}

	url := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPISearchBrandkit
	utils.Log(fmt.Sprintf("Count Brandkit Request: Url : %s  , request : %s  ,headers : %+v ", url, string(jsonData), headers))
	response, err := sendRequest(url, string(jsonData), headers)
	if err != nil {
		return brandKitData, err
	}

	var brandKitAPIResponse brandkit_model.BrandKitAPIResponse
	if err := json.Unmarshal([]byte(response), &brandKitAPIResponse); err != nil {
		return brandKitData, err
	}

	count, err := get_count(netcoreApiRequestBody, url, headers)

	if err != nil {
		return brandkit_model.BrandKitData{}, err
	}

	data := brandkit_model.BrandKitData{
		Count: &count.Count,
	}

	return data, nil
}
