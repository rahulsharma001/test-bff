package fonts_model

type Fonts struct {
	Id         int64  `json:"id"`
	FontName   string `json:"font_name"`
	FontFamily string `json:"font_family"`
	Url        string `json:"font_url"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

type FontsAPIResponse struct {
	RequestID string  `json:"request_id"`
	Status    string  `json:"status"`
	Message   string  `json:"message"`
	Code      int     `json:"code"`
	Fonts     []Fonts `json:"data,omitempty"`
	Count     int     `json:"count,omitempty"`
	Errors    string  `json:"errors,omitempty"`
}

type GetList struct {
	Page      int64  `form:"page" binding:"gte=0"`
	Limit     int64  `form:"limit" binding:"gte=0"`
	SearchFor string `form:"search_for,omitempty"`
}

func NewGetListRequest() *GetList {
	return &GetList{
		Page:      0,
		Limit:     10,
		SearchFor: "",
	}
}
