package auth_model

type DataMaskSettings struct {
	Fields  int `json:"fields"`  // Email=2, Mobile=4, Both=1, None=0
	Options int `json:"options"` // Ui=2, Download=4, Both=1, None=0
}

type Settings struct {
	DataMask DataMaskSettings `json:"data_mask"`
}

type ResponseData struct {
	ID              int      `json:"id"`
	Username        string   `json:"username"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	APIKey          string   `json:"api_key"`
	ReferenceUserID int      `json:"reference_user_id"`
	Settings        Settings `json:"settings"`
}

type Response struct {
	Status int          `json:"status"`
	Data   ResponseData `json:"data"`
}
