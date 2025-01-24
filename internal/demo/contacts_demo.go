package demo

import "github.com/gin-gonic/gin"

// RegisterContactsDemoData registers demo data for the /contacts/get endpoint
func RegisterContactsDemoData(c *gin.Context) {
	DemoResponses["/v1/contacts/get_contacts"] = `{
    "code": 200,
    "data": {
        "count": "3543108",
        "results": [
            {
                "blacklisted": "0",
                "blacklisted_mobile": "0",
                "email": null,
                "entered": "2024-09-05 15:35:38",
                "foreignkey": "04BCB9D5-8358-4A4D-B2A7-D23B6AFF4789",
                "id": "99420148",
                "mobile": "7878787877",
                "modified": "2024-09-05 15:35:38"
            },
            {
                "blacklisted": "0",
                "blacklisted_mobile": "0",
                "email": "jshaasd@jashdas.sajiuh",
                "entered": "2024-09-13 16:06:12",
                "foreignkey": "jshaasd@jashdas.sajiuh",
                "htmlemail": "1",
                "id": "13965124",
                "mobile": "9865421273",
                "modified": "2024-09-13 16:06:12",
                "passwordchanged": "2024-09-13",
                "uniqid": "40c7383a982fc6bbf4c82efaa1090aa7",
                "widgetflag": "11"
            }
        ],
        "tableHeaders": [
            {
                "attribute_name": "GUID",
                "attribute_value": "GUID",
                "class": "text-left",
                "data_type": "varchar",
                "data_type_with_limit": "varchar(255)",
                "default_value": null,
                "id": -1,
                "is_foreignkey": true,
                "is_special": true,
                "key": "foreignkey",
                "order": null,
                "special_order": 1,
                "text": "GUID",
                "type": "link"
            },
            {
                "attribute_name": "mobile",
                "attribute_value": "MOBILE",
                "class": "text-left",
                "data_type": "varchar",
                "data_type_with_limit": "varchar(255)",
                "default_value": null,
                "id": -2,
                "is_foreignkey": false,
                "is_special": true,
                "key": "mobile",
                "order": null,
                "special_order": 3,
                "text": "MOBILE",
                "type": "data"
            }
        ]
    },
    "message": "data fetched successfully.",
    "status": "success"
}`

}
