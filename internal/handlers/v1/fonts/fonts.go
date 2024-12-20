package font_handler

import (
	fonts_model "cee-bff-go/internal/models/v1/fonts"
	fonts_service "cee-bff-go/internal/services/v1/fonts"
	"cee-bff-go/internal/utils"
	"cee-bff-go/internal/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setResponse(ctx *gin.Context, statusCode int, status, message string, data interface{}) {
	utils.SetResponse(ctx, utils.CommonResponse{
		StatusCode: statusCode,
		Status:     status,
		Message:    message,
		Data:       data,
	})
}

func Create(ctx *gin.Context) {

	var fonts []fonts_model.Fonts

	err := ctx.ShouldBindJSON(&fonts)

	if err != nil {
		utils.Log("Error while searching font: " + err.Error())
		utils.SetResponse(ctx, utils.CommonResponse{
			StatusCode: http.StatusAlreadyReported,
			Status:     "error",
			Message:    "Failed to create font",
			Data:       nil,
		})
		return
	}

	err = fonts_service.Create(fonts)

	if err != nil {
		utils.Log("Error while searching font: " + err.Error())
		utils.SetResponse(ctx, utils.CommonResponse{
			StatusCode: http.StatusAlreadyReported,
			Status:     "error",
			Message:    "Failed to create font",
			Data:       nil,
		})
		return
	}

	utils.SetResponse(ctx, utils.CommonResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "Font created successfully.",
		Data:       nil,
	})
}

func GetList(ctx *gin.Context) {
	request := fonts_model.NewGetListRequest()
	v := validators.NewValidator()

	if err := fonts_service.ValidateGetList(v, ctx, request); err != nil {
		utils.Log("Error validating GetList request: " + err.Error())
		setResponse(ctx, http.StatusBadRequest, "error", "Validation failed", nil)
		return
	}

	allFonts, err := fonts_service.GetAll(request)

	if err != nil {
		utils.Log("Error while searching font: " + err.Error())
		utils.SetResponse(ctx, utils.CommonResponse{
			StatusCode: http.StatusAlreadyReported,
			Status:     "error",
			Message:    "Failed to fetch fonts",
			Data:       nil,
		})
		return
	}

	font_count := len(allFonts)
	utils.SetResponse(ctx, utils.CommonResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "Fonts fetched successfully",
		Data:       allFonts,
		Count:      &font_count,
	})
}
