package netcoreapi_model

import brandkit_model "cee-bff-go/internal/models/v1/brandkit"

type ConditionDetail struct {
	Field         string      `json:"field"`
	FieldCategory string      `json:"field_category"`
	Operation     string      `json:"operation"`
	Value         interface{} `json:"value"`
}

// Struct for filtering criteria
type FilteringCriteria struct {
	ConditionOperator string            `json:"condition_operator"`
	ConditionDetails  []ConditionDetail `json:"condition_details"`
}

// Struct for sorting options
type SortingOption struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// Struct for pagination options
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// Struct for output configuration
type OutputConfig struct {
	GetCount   bool            `json:"get_count"`
	Sorting    []SortingOption `json:"sorting,omitempty"`
	Pagination *Pagination     `json:"pagination,omitempty"`
	Fields     []string        `json:"fields"`
}

type Data struct {
	brandkit_model.BrandKit
}

// Struct for the common request
type APIRequest struct {
	FilteringCriteriaOperator string              `json:"filtering_criteria_operator,omitempty"`
	FilteringCriteria         []FilteringCriteria `json:"filtering_criteria,omitempty"`
	Output                    *OutputConfig       `json:"output,omitempty"`
	Data                      interface{}         `json:"data,omitempty"`
}
