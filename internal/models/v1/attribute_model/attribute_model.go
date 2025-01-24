package attribute_model

type UpdateAttributeOrderRequest struct {
	PageURI     string         `json:"page_uri"`
	PageConfigs AttributeSlice `json:"page_configs"`
}

type AttributeOrderResponse struct {
	Code    int            `json:"code"`
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    AttributeSlice `json:"data"`
}

type InfoTooltip struct {
	Text      string `json:"text"`
	Class     string `json:"class"`
	Placement string `json:"placement"`
}

type Attribute struct {
	ID                  any         `json:"id"`
	AttributeName       string      `json:"attribute_name"`
	AttributeValue      string      `json:"attribute_value"`
	AttributeType       string      `json:"attribute_type"`
	Type                string      `json:"type"`
	SystemAttributeName string      `json:"system_attribute_name"`
	DefaultValue        interface{} `json:"default_value"`
	Segmentation        string      `json:"segmentation"`
	Value               string      `json:"value"`
	Name                string      `json:"name"`
	AttributeCategory   string      `json:"attribute_category"`
	Hideable            bool        `json:"hideable"`
	Visible             bool        `json:"visible"`
	SubText             interface{} `json:"subText"`
	Order               *int        `json:"order"`
	Code                string      `json:"code"`
	IsSpecial           bool        `json:"is_special"`
	IsSelected          interface{} `json:"is_selected"` // Can be bool or numeric, use interface{}
	IsForeignKey        bool        `json:"is_foreignkey"`
	SpecialOrder        *int        `json:"special_order"`
}

type AttributeSlice []Attribute

func (attributes AttributeSlice) GetAttributesNames() []string {
	var selectedAttributes []string
	for _, attribute := range attributes {
		if attribute.Type != "longtext" {
			selectedAttributes = append(selectedAttributes, attribute.AttributeValue)
		}
	}
	return selectedAttributes
}

func (attributes AttributeSlice) GetSelectedAttributesNames() []string {
	var selectedAttributes []string
	for _, attribute := range attributes {
		if attribute.Visible {
			selectedAttributes = append(selectedAttributes, attribute.AttributeValue)
		}
	}
	return selectedAttributes
}

func (attributes AttributeSlice) GetSelectedAttributes() AttributeSlice {
	var selectedAttributes AttributeSlice
	for _, attribute := range attributes {
		if attribute.Visible {
			selectedAttributes = append(selectedAttributes, attribute)
		}
	}
	return selectedAttributes
}

func (attributes AttributeSlice) GetAttributesWithoutLongText() AttributeSlice {
	var selectedAttributes AttributeSlice
	for _, attribute := range attributes {
		if attribute.Type != "longtext" {
			selectedAttributes = append(selectedAttributes, attribute)
		}
	}
	return selectedAttributes
}

func (attributes AttributeSlice) GetAttributesWithTableHeaders() []TableHeaderAttribute {
	var tableHeaders []TableHeaderAttribute

	for _, attribute := range attributes {
		// Determine the type based on the attribute name
		var fieldType string
		var infotooltip *InfoTooltip
		if attribute.Name == "GUID" {
			fieldType = "link"
			infotooltip = &InfoTooltip{
				Text:      "Global Unique Identifier (GUID) is a random ID assigned to an App/Web user. GUID is not assigned for manual upload when contacts do not have the primary key.",
				Class:     "nxt-white-tooltip nxt-anon-guid-tooltip",
				Placement: "right",
			}
		} else {
			fieldType = "data"
		}

		// Create a new TableHeaderAttribute and populate it based on the given rules
		header := TableHeaderAttribute{
			AttributeName:     attribute.AttributeValue,
			AttributeValue:    attribute.AttributeValue,
			Class:             "text-left",             // Hardcoded
			DataType:          attribute.Type,          // Taken from type field
			DataTypeWithLimit: attribute.AttributeType, // Taken from attribute_type field
			DefaultValue:      attribute.DefaultValue,
			ID:                attribute.ID, // Assuming ID can be any type, helper function used
			IsForeignKey:      attribute.IsForeignKey,
			IsSpecial:         attribute.IsSpecial,
			Key:               attribute.AttributeValue,
			Order:             attribute.Order,
			SpecialOrder:      attribute.SpecialOrder,   // Handle nil pointer for special_order
			Text:              attribute.AttributeValue, // Assuming this maps to 'Value' field
			Type:              fieldType,                // Determined based on the name
			InfoTooltip:       infotooltip,
		}

		// Append the populated TableHeaderAttribute to the slice
		tableHeaders = append(tableHeaders, header)
	}

	return tableHeaders
}

type TableHeaderAttribute struct {
	AttributeName     string       `json:"attribute_name"`
	AttributeValue    string       `json:"attribute_value"`
	Class             string       `json:"class"`
	DataType          string       `json:"data_type"`
	DataTypeWithLimit string       `json:"data_type_with_limit"`
	DefaultValue      interface{}  `json:"default_value"` // Use `interface{}` for null or dynamic values
	ID                any          `json:"id"`
	IsForeignKey      bool         `json:"is_foreignkey"`
	IsSpecial         bool         `json:"is_special"`
	Key               string       `json:"key"`
	Order             interface{}  `json:"order"` // Use `interface{}` for null
	SpecialOrder      *int         `json:"special_order"`
	SubText           string       `json:"sub_text"`
	Text              string       `json:"text"`
	Type              string       `json:"type"`
	InfoTooltip       *InfoTooltip `json:"infoTooltip"`
}

type UserAttributePayload struct {
	PageURI     string `json:"page_uri"`
	ContactType string `json:"contact_type"`
}
