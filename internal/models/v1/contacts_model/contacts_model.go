package contacts_model

import (
	common_model "cee-bff-go/internal/models"
	"errors"
	"slices"
	"strconv"
	"time"
)

type GetContactsRequest struct {
	Limit                      int      `json:"limit"`
	Offset                     int      `json:"offset"`
	Page                       int      `json:"page"`
	Type                       string   `json:"type"`
	Field                      string   `json:"field"`
	Fields                     []string `json:"fields"`
	Direction                  string   `json:"direction"`
	SearchFor                  string   `json:"search_for"`
	SegmentID                  int      `json:"list_segment_type_drop_down,omitempty"`
	AttributeDropDown          string   `json:"attribute_drop_down,omitempty"`
	AttributeConditionDropDown string   `json:"attribute_condition_drop_down,omitempty"`
	AttributeConditionValue    string   `json:"attribute_condition_value,omitempty"`
	Blacklist                  bool     `json:"blacklist,omitempty"`
	DateCondition              string   `json:"date_condition,omitempty"`
	DateStatus                 string   `json:"date_status,omitempty"`
	LastDayValue               int      `json:"last_day_value,omitempty"`
	FromDate                   string   `json:"frmdate,omitempty"`
	ToDate                     string   `json:"todate,omitempty"`
}

func (gcr GetContactsRequest) GetSearchCondition() ([]common_model.ConditionDetail, error) {
	var searchCondition []common_model.ConditionDetail

	// Searching condition for attribute search
	if gcr.AttributeConditionDropDown != "" || gcr.AttributeConditionValue != "" || gcr.AttributeDropDown != "" {
		categoryMap := []string{
			"GUID",
			"MOBILE",
			"EMAIL",
			"CUSTOMER_ID",
		}

		var category string
		exists := slices.Contains(categoryMap, gcr.AttributeDropDown)
		if exists {
			category = "config"
		} else {
			category = "attribute"
		}

		operation := map[string]string{
			"is":          "equals",
			"like":        "contains",
			"begins_with": "starts_with",
			"ends_with":   "ends_with",
		}

		searchCondition = append(searchCondition, common_model.ConditionDetail{
			Field:         gcr.AttributeDropDown,
			FieldCategory: category,
			Operation:     operation[gcr.AttributeConditionDropDown],
			Value:         []string{gcr.AttributeConditionValue},
		})
	}

	// Adding condition for segment_id
	if gcr.SegmentID != 0 {
		searchCondition = append(searchCondition, common_model.ConditionDetail{
			Field:         "segment_id",
			FieldCategory: "audience",
			Operation:     "equals",
			Value:         []string{strconv.Itoa(gcr.SegmentID)},
		})
	}

	// Adding conditions based on date_condition
	if gcr.DateCondition == "ever" {
		searchCondition = append(searchCondition, common_model.ConditionDetail{
			Field:         "created_at",
			FieldCategory: "config",
			Operation:     "is_not_null",
			Value:         []string{},
		})
	} else if gcr.DateCondition == "specific_date" || gcr.DateCondition == "date_range" {
		// Make sure `FromDate` and `ToDate` are properly formatted
		if gcr.FromDate != "" && gcr.ToDate != "" {
			searchCondition = append(searchCondition, common_model.ConditionDetail{
				Field:         "created_at",
				FieldCategory: "config",
				Operation:     "between",
				Value:         []string{gcr.FromDate + " 00:00:00", gcr.ToDate + " 23:59:59"},
			})
		}
	} else if gcr.DateCondition == "in_last" {
		// `LastDayValue` is now an integer
		if gcr.LastDayValue > 0 {
			// Get today's date
			today := time.Now()

			// Calculate the date for "last X days"
			startDate := today.AddDate(0, 0, -gcr.LastDayValue).Format("2006-01-02")

			// Add condition for "in_last" date range
			searchCondition = append(searchCondition, common_model.ConditionDetail{
				Field:         "created_at",
				FieldCategory: "config",
				Operation:     "between",
				Value:         []string{startDate + " 00:00:00", today.Format("2006-01-02") + " 23:59:59"},
			})
		}
	}

	return searchCondition, nil
}

func (gcr GetContactsRequest) GetFilteringCriteria() ([]common_model.FilteringCriteria, error) {
	var filteringCriteria []common_model.FilteringCriteria
	var conditionDetails []common_model.ConditionDetail

	// 1. Get the search conditions (e.g., segment_id, created_at, etc.)
	searchConditions, err := gcr.GetSearchCondition()
	if err != nil {
		return nil, err
	}

	// 2. Add "contact_type" condition if present
	if gcr.Type != "" {
		contactTypeCondition := common_model.ConditionDetail{
			Field:         "contact_type",
			FieldCategory: "config",
			Operation:     "equals",
			Value:         []string{gcr.Type},
		}
		conditionDetails = append(conditionDetails, contactTypeCondition)
	}

	// 3. Add the conditions obtained from GetSearchCondition()
	if len(searchConditions) > 0 {
		conditionDetails = append(conditionDetails, searchConditions...)
	}

	// 4. Add a separate filtering group for SearchFor (if provided)
	conditionOperator := "and"
	if gcr.SearchFor != "" {
		conditionOperator = "or"
		searchForCondition := []common_model.ConditionDetail{
			{
				Field:         "email",
				FieldCategory: "config",
				Operation:     "contains",
				Value:         []string{gcr.SearchFor},
			},
			{
				Field:         "guid",
				FieldCategory: "config",
				Operation:     "contains",
				Value:         []string{gcr.SearchFor},
			},
			{
				Field:         "mobile",
				FieldCategory: "config",
				Operation:     "contains",
				Value:         []string{gcr.SearchFor},
			},
		}

		filteringCriteria = append(filteringCriteria, common_model.FilteringCriteria{
			ConditionOperator: conditionOperator,
			ConditionDetails:  searchForCondition,
		})
	}

	// If there are any general conditions, append them to the filtering criteria
	if len(conditionDetails) > 0 {
		filteringCriteria = append(filteringCriteria, common_model.FilteringCriteria{
			ConditionOperator: conditionOperator,
			ConditionDetails:  conditionDetails,
		})
	}

	// If no filtering criteria were added, return an error
	if len(filteringCriteria) == 0 {
		return nil, errors.New("no filtering criteria provided")
	}

	return filteringCriteria, nil
}

type GetContactAPIResponse struct {
	Count       int         `json:"count"`
	Results     interface{} `json:"results"`
	TableHeader interface{} `json:"tableHeaders"`
}

type NetcoreContactRequest struct {
	FilteringCriteriaOperator string                           `json:"filtering_criteria_operator"`
	FilteringCriteria         []common_model.FilteringCriteria `json:"filtering_criteria"`
	Output                    common_model.OutputFields        `json:"output,omitempty"`
}

// Define a struct to hold the response from the Netcore API
type NetcoreContactResponse struct {
	RequestID string                   `json:"request_id"`
	Code      int                      `json:"code"`
	Status    string                   `json:"status"`
	Data      []map[string]interface{} `json:"data"`
	Count     int                      `json:"count,omitempty"`
}

// FilterNAValues filters out keys with the value "NA" from all maps in Data.
func (ncr *NetcoreContactResponse) FilterNAValues() map[string]string {
	if len(ncr.Data) == 0 {
		return nil
	}

	// Initialize the result map
	result := make(map[string]string)

	// Iterate over all maps in the Data slice
	for _, m := range ncr.Data {
		for key, value := range m {
			// Ensure that the value is a string before comparing with "NA"
			if strValue, ok := value.(string); ok && strValue != "NA" {
				result[key] = strValue
			}
		}
	}

	return result
}

// ReplaceGUIDNAWithNotApplicable checks if the key "GUID" has the value "NA" and replaces it with "Not Applicable"
func (ncr *NetcoreContactResponse) ReplaceGUIDValues() {
	// Iterate over each map in the Data slice
	for _, m := range ncr.Data {
		// Check if the "GUID" key exists and if its value is "NA"
		if value, exists := m["GUID"]; exists {
			if strValue, ok := value.(string); ok && strValue == "NA" {
				// Replace "NA" with "Not Applicable"
				m["GUID"] = "GUID not assigned"
			}
		}
	}
}

type NetcoreAudienceSearchResponse struct {
	RequestID string     `json:"request_id"`
	Code      int        `json:"code"`
	Status    string     `json:"status"`
	Data      []Audience `json:"data"`
	Count     int        `json:"count,omitempty"`
}

type Audience struct {
	AudienceID   int    `json:"audience_id,omitempty"`
	AudienceName string `json:"audience_name,omitempty"`
}

func (n NetcoreAudienceSearchResponse) GetAudienceSegmentMap() map[int]string {
	var audienceMap map[int]string = make(map[int]string)
	for _, v := range n.Data {
		audienceMap[v.AudienceID] = v.AudienceName
	}
	return audienceMap
}
