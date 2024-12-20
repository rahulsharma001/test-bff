package common_model

type FilterRequest struct {
	AudienceType              string              `json:"audience_type,omitempty"`
	FilteringCriteriaOperator string              `json:"filtering_criteria_operator"`
	FilteringCriteria         []FilteringCriteria `json:"filtering_criteria"`
	Output                    OutputFields        `json:"output,omitempty"`
}

type FilteringCriteria struct {
	ConditionOperator string            `json:"condition_operator"`
	ConditionDetails  []ConditionDetail `json:"condition_details"`
}

type ConditionDetail struct {
	Field         string   `json:"field"`
	FieldCategory string   `json:"field_category"`
	Operation     string   `json:"operation"`
	Value         []string `json:"value"`
}

type OutputFields struct {
	GetCount   bool        `json:"get_count,omitempty"`
	Sorting    []*Sorting  `json:"sorting,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Fields     []string    `json:"fields,omitempty"`
}

type Pagination struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

type Sorting struct {
	Field     string `json:"field,omitempty"`
	Direction string `json:"direction,omitempty"`
}

type NetcoreAPIResponse struct {
	RequestID string      `json:"request_id"`
	Code      int         `json:"code"`
	Status    string      `json:"status"`
	Data      interface{} `json:"data"`
	Count     int         `json:"count,omitempty"`
}
