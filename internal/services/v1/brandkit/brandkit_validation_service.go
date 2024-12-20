package service

import (
	brandkit_model "cee-bff-go/internal/models/v1/brandkit"
	"cee-bff-go/internal/utils"
	"cee-bff-go/internal/validators"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateGetList(v *validators.Validator, c *gin.Context, request *brandkit_model.GetList) error {
	err := c.ShouldBindJSON(&request)

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

func ValidateGet(v *validators.Validator, c *gin.Context, request *brandkit_model.Get) error {

	if err := c.ShouldBindUri(&request); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, validators.CustomErrorMessage(err))
		}
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
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Validation failed"})
		return errors.New("validation failed")
	}
	return nil
}
