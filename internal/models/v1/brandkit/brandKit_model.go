package brandkit_model

type Font struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type HeadingFont struct {
	FontName     map[string]interface{} `json:"fontName"`
	FontSize     string                 `json:"fontSize"`
	FontWeight   Font                   `json:"fontWeight"`
	FallbackFont map[string]interface{} `json:"fallbackFont"`
}

type ParagraphFonts struct {
	FontName     map[string]interface{} `json:"fontName"`
	FontSize     string                 `json:"fontSize"`
	FontWeight   Font                   `json:"fontWeight"`
	FallbackFont map[string]interface{} `json:"fallbackFont"`
}

type Fonts struct {
	HeadingFonts   map[string]HeadingFont `json:"headingFonts"`
	ParagraphFonts ParagraphFonts         `json:"paragraphFonts"`
}

type Logos struct {
	HeaderLogo struct {
		Urls []string `json:"urls"`
	} `json:"headerLogo"`
	FooterLogo struct {
		UseSameAsHeader bool     `json:"useSameAsHeader"`
		Urls            []string `json:"urls"`
	} `json:"footerLogo"`
}

type Colors struct {
	ColorPalette  []string `json:"colorPalette"`
	BodyTextColor []string `json:"bodyTextColor"`
	LinksColor    []string `json:"linksColor"`
}

type SocialLink struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Icon   string `json:"icon"`
	Active bool   `json:"active"`
	Custom bool   `json:"custom"`
}

type ButtonSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
type SocialLinks struct {
	Platforms   []SocialLink `json:"platforms"`
	ButtonSize  string       `json:"buttonSize"`
	ButtonStyle string       `json:"buttonStyle"`
}

type HeaderLink struct {
	Text string `json:"text"`
	Link string `json:"link"`
}

type FooterText struct {
	Text           string `json:"text"`
	AdditionalInfo string `json:"additionalInfo"`
}

type Details struct {
	Logos       Logos                  `json:"logos"`
	Colors      Colors                 `json:"colors"`
	Fonts       Fonts                  `json:"fonts"`
	Buttons     map[string]interface{} `json:"buttons"`
	SocialLinks SocialLinks            `json:"socialLinks"`
	HeaderLinks []HeaderLink           `json:"headerLinks"`
	FooterText  []FooterText           `json:"footerText"`
}

type BrandKitId struct {
	Id int `json:"id,omitempty"`
}

type BrandKit struct {
	BrandKitId
	Name      string  `json:"name" binding:"required"`
	Status    *bool   `json:"status,omitempty"`
	Details   Details `json:"details,omitempty" binding:"required"`
	CreatedAt string  `json:"created_at,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}
type BrandKitResponse struct {
	BrandKitId
	Name      string  `json:"name" binding:"required"`
	Status    *bool   `json:"status,omitempty"`
	Details   Details `json:"details,omitempty" binding:"required"`
	CreatedAt string  `json:"created_at,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

type BrandKitData struct {
	Data  []BrandKit `json:"kit,omitempty"`
	Count *int       `json:"count,omitempty"`
}

type BrandKitAPIResponse struct {
	RequestID string     `json:"request_id"`
	Status    string     `json:"status"`
	Message   string     `json:"message"`
	Code      int        `json:"code"`
	BrandKits []BrandKit `json:"data,omitempty"`
	Count     int        `json:"count,omitempty"`
	Errors    string     `json:"errors,omitempty"`
}
type NetcoreAPIBrandKitResponse struct {
	RequestID string               `json:"request_id"`
	Status    string               `json:"status"`
	Message   string               `json:"message"`
	Code      int                  `json:"code"`
	BrandKits []NetcoreAPIBrandKit `json:"data,omitempty"`
	Count     int                  `json:"count,omitempty"`
	Errors    string               `json:"errors,omitempty"`
}

type NetcoreAPIBrandKit struct {
	BrandKitId
	Name      string `json:"name" binding:"required"`
	Status    *bool  `json:"status,omitempty"`
	Details   string `json:"details,omitempty" binding:"required"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type SearchRequest struct {
	SearchTerm string `form:"search_for"`
}

type GetList struct {
	Page      int64  `json:"page" binding:"gte=0"`
	Limit     int64  `json:"limit" binding:"gte=0"`
	SearchFor string `json:"search_for,omitempty"`
}

type Get struct {
	Id int64 `uri:"id" binding:"gt=0"`
}

type Create struct {
	Name    string  `json:"name" binding:"required"`
	Details Details `json:"details" binding:"required"`
}
