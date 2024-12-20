package service

import (
	font_model "cee-bff-go/internal/models/v1/fonts"
	"cee-bff-go/internal/utils"
	"cee-bff-go/internal/validators"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateGetList(v *validators.Validator, c *gin.Context, request *font_model.GetList) error {
	err := c.ShouldBindQuery(&request)
	if err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, validators.CustomErrorMessage(err))
		}
		//c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "Invalid request",
			Error:      validationErrors,
		})
		return errors.New("invalid request")
	}

	if err := v.Validate(request); err != nil {
		utils.SetResponse(c, utils.CommonResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "Validation failed",
		})
		return errors.New("validation failed")
	}
	return nil
}
