package service

import (
	"cee-bff-go/internal/config/endpoints"
	fonts_model "cee-bff-go/internal/models/v1/fonts"
	netcoreapi_model "cee-bff-go/internal/models/v1/netcoreapi"
	"cee-bff-go/internal/utils"
	"encoding/json"
	"time"

	"github.com/spf13/viper"
)

func prepareHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"api-key":      viper.GetString("NETCORE_API_KEY"),
	}
}

func Create(fonts []fonts_model.Fonts) error {
	headers := prepareHeaders()
	netcoreApiRequestBody := netcoreapi_model.APIRequest{Data: fonts}
	byteData, _ := json.Marshal(netcoreApiRequestBody)
	url := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPICreateFont
	httpParams := utils.NewHTTP("POST", url, headers, string(byteData))
	httpParams.Timeout = 10 * time.Second
	response, err := httpParams.DoHTTP()
	if err != nil {
		return err
	}
	var fontAPIResponse fonts_model.FontsAPIResponse
	err = json.Unmarshal([]byte(response), &fontAPIResponse)
	if err != nil {
		return err
	}
	return nil
}

func GetAll(request *fonts_model.GetList) ([]map[string]interface{}, error) {

	headers := prepareHeaders()
	conditionDetails := []netcoreapi_model.ConditionDetail{}
	if request.SearchFor != "" {
		conditionDetails = append(conditionDetails, netcoreapi_model.ConditionDetail{
			Field:         "font_name",
			FieldCategory: "config",
			Operation:     "like",
			Value:         []interface{}{request.SearchFor},
		}, netcoreapi_model.ConditionDetail{
			Field:         "font_family",
			FieldCategory: "config",
			Operation:     "like",
			Value:         []interface{}{request.SearchFor},
		})
	}
	netcoreApiRequestBody := createRequestBody("or", "or", conditionDetails, false, int(request.Page), int(request.Limit), true)

	byteData, err := json.Marshal(netcoreApiRequestBody)

	if err != nil {
		return nil, err
	}
	url := viper.GetString("NETCORE_API_ENDPOINT") + endpoints.NetcoreAPISearchFont
	httpParams := utils.NewHTTP("POST", url, headers, string(byteData))

	response, err := httpParams.DoHTTP()
	if err != nil {
		return nil, err
	}

	var fontAPIResponse fonts_model.FontsAPIResponse
	err = json.Unmarshal([]byte(response), &fontAPIResponse)
	if err != nil {
		return nil, err
	}
	var fontResp []map[string]interface{}
	for _, data := range fontAPIResponse.Fonts {
		fontResult := make(map[string]interface{})
		fontResult["id"] = data.Id
		fontResult["font_name"] = data.FontName
		fontResult["name"] = data.FontName
		fontResult["font_family"] = data.FontFamily
		fontResult["font_url"] = data.Url
		fontResult["created_at"] = data.CreatedAt
		fontResult["updated_at"] = data.UpdatedAt
		fontResp = append(fontResp, fontResult)
	}

	return fontResp, nil
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
					Field:     "updated_at",
					Direction: "desc",
				},
			},
			Pagination: &netcoreapi_model.Pagination{
				Page:  page,
				Limit: limit,
			},
			Fields: []string{"id", "font_family", "font_name", "font_url", "created_at", "updated_at"},
		}
	}
	return request
}
