package demo

import (
	"github.com/gin-gonic/gin"
)

func RegisterAttributesDemoData(c *gin.Context) {
	DemoResponses["/v1/attributes/get_attributes"] = `{
    "code": 200,
    "status": "success",
    "message": "Fetched Attributes Successfully.",
    "data": [
        {
            "id": -3,
            "attribute_name": "guid",
            "attribute_value": "GUID",
            "attribute_type": "varchar(255)",
            "type": "varchar",
            "system_attribute_name": "GUID",
            "default_value": null,
            "segmentation": "1",
            "value": "guid",
            "name": "GUID",
            "attribute_category": "system",
            "hideable": false,
            "visible": true,
            "subText": null,
            "order": -1,
            "code": "guid",
            "is_special": true,
            "is_selected": true,
            "is_foreignkey": true,
            "special_order": 1
        },
        {
            "id": -2,
            "attribute_name": "mobile",
            "attribute_value": "MOBILE",
            "attribute_type": "varchar(255)",
            "type": "varchar",
            "system_attribute_name": "MOBILE",
            "default_value": null,
            "segmentation": "1",
            "value": "mobile",
            "name": "MOBILE",
            "attribute_category": "system",
            "hideable": false,
            "visible": true,
            "subText": null,
            "order": -2,
            "code": "mobile",
            "is_special": true,
            "is_selected": true,
            "is_foreignkey": false,
            "special_order": 3
        },
        {
            "id": 1,
            "attribute_name": "DATEOFA",
            "attribute_value": "DATEOFA",
            "attribute_type": "varchar(255)",
            "type": "varchar",
            "system_attribute_name": null,
            "default_value": null,
            "segmentation": "1",
            "value": "DATEOFA",
            "name": "DATEOFA",
            "attribute_category": "custom",
            "hideable": false,
            "visible": true,
            "subText": null,
            "order": 1,
            "code": "DATEOFA",
            "is_special": true,
            "is_selected": true,
            "is_foreignkey": false,
            "special_order": 0
        }
    ]
}`
}
