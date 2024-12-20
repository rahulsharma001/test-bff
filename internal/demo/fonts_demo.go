package demo

import (
	"github.com/gin-gonic/gin"
)

// RegisterContactsDemoData registers demo data for the /contacts/get endpoint
func RegisterFontsDemoData(c *gin.Context) {
	DemoResponses["/v1/fonts/"] = `{"request_id":"2a78eec8-4070-4d30-83ca-f5f4a78f8668","code":200,"status":"success","message":"Fonts fetched successfully","data":[{"created_at":"2024-11-21 09:57:19","font_family":"sans","font_name":"Open sans","id":1,"name":"Open sans","updated_at":"2024-11-21 09:57:19","font_url":"https://fonts.gstatic.com/s/opensans/v34/memSYaGs126MiZpBA-UvWbX2vVnXBbObj2OVZyOOSr4dVJWUgsjZ0B4gaVI.woff2"},{"created_at":"2024-11-21 09:57:19","font_family":"sans","font_name":"Open sans","id":1,"name":"Open sans","updated_at":"2024-11-21 09:57:19","font_url":"https://fonts.gstatic.com/s/opensans/v34/memSYaGs126MiZpBA-UvWbX2vVnXBbObj2OVZyOOSr4dVJWUgsjZ0B4gaVI.woff2"},{"created_at":"2024-11-21 09:57:19","font_family":"sans","font_name":"Open sans","id":1,"name":"Open sans","updated_at":"2024-11-21 09:57:19","font_url":"https://fonts.gstatic.com/s/opensans/v34/memSYaGs126MiZpBA-UvWbX2vVnXBbObj2OVZyOOSr4dVJWUgsjZ0B4gaVI.woff2"}]}`

	DemoResponses["/v1/fonts/create"] = `{"code":200,"status":"success","message":"Font created successfully."}`

}
