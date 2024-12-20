package brandkit

import (
	brandkit_model "cee-bff-go/internal/models/v1/brandkit"
	brandkit_service "cee-bff-go/internal/services/v1/brandkit"
	"cee-bff-go/internal/utils"
	"cee-bff-go/internal/validators"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Helper function to set the response with a common format
func setResponse(ctx *gin.Context, statusCode int, status, message string, data interface{}) {
	utils.SetResponse(ctx, utils.CommonResponse{
		StatusCode: statusCode,
		Status:     status,
		Message:    message,
		Data:       data,
	})
}

func GetList(ctx *gin.Context) {
	var request brandkit_model.GetList
	v := validators.NewValidator()

	if err := brandkit_service.ValidateGetList(v, ctx, &request); err != nil {
		utils.Log("Error validating GetList request: " + err.Error())
		setResponse(ctx, http.StatusBadRequest, "error", "Validation failed", nil)
		return
	}

	allBrandkits, err := brandkit_service.GetAll(&request)
	if err != nil {
		utils.Log("Failed to fetch brand kits: " + err.Error())
		setResponse(ctx, http.StatusBadRequest, "error", "Failed to fetch brand kits", nil)
		return
	}

	setResponse(ctx, http.StatusOK, "success", "Brand kits fetched successfully", allBrandkits)
}

func Get(ctx *gin.Context) {
	var request brandkit_model.Get
	v := validators.NewValidator()

	if err := brandkit_service.ValidateGet(v, ctx, &request); err != nil {
		utils.Log("Failed to validate request: " + err.Error())
		setResponse(ctx, http.StatusBadRequest, "error", "Invalid request", nil)
		return
	}

	brandkit, err := brandkit_service.Get(request.Id)
	if err != nil {
		utils.Log("Failed to fetch brand kit: " + err.Error())
		setResponse(ctx, http.StatusNotAcceptable, "error", "Failed to fetch brand kit", nil)
		return
	}

	setResponse(ctx, http.StatusOK, "success", "Brand kit fetched successfully", brandkit)
}

func GetActive(ctx *gin.Context) {
	brandkit, err := brandkit_service.GetActive()
	if err != nil {
		setResponse(ctx, http.StatusNotAcceptable, "error", "Failed to fetch active brand kit", nil)
		return
	}

	setResponse(ctx, http.StatusOK, "success", "Active brand kit fetched successfully", brandkit)
}

func Create(ctx *gin.Context) {
	var brandKit brandkit_service.BrandKit
	err := ctx.ShouldBindJSON(&brandKit)
	fmt.Println("Errors ==>", err)
	if err != nil {
		utils.Log("Error parsing Create request: " + err.Error())
		setResponse(ctx, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	if err := brandKit.Create(); err != nil {
		fmt.Println("Errors ==>", err)
		utils.Log("Error while creating brandkit: " + err.Error())
		setResponse(ctx, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	setResponse(ctx, http.StatusOK, "success", "Brand kit created successfully", nil)
}

func Edit(ctx *gin.Context) {
	var brandKit brandkit_service.BrandKit
	if err := ctx.ShouldBindJSON(&brandKit); err != nil {
		utils.Log("Error while parsing brandkit: " + err.Error())
		setResponse(ctx, http.StatusAlreadyReported, "error", "Failed to update brand kit", nil)
		return
	}

	if err := brandKit.Edit(); err != nil {
		utils.Log("Error while editing brandkit: " + err.Error())
		setResponse(ctx, http.StatusAlreadyReported, "error", "Failed to update brand kit", nil)
		return
	}

	setResponse(ctx, http.StatusOK, "success", "Brand kit updated successfully", nil)
}

func Delete(ctx *gin.Context) {
	kitid, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Log("Error while parsing brandkit: " + err.Error())
		setResponse(ctx, http.StatusNotAcceptable, "error", "Invalid brand kit ID", nil)
		return
	}

	err = brandkit_service.Delete(kitid)
	fmt.Println(err)
	if err != nil {
		utils.Log("Error while deleting brandkit: " + err.Error())
		setResponse(ctx, http.StatusNotAcceptable, "error", "Failed to delete brand kit", nil)
		return
	}

	setResponse(ctx, http.StatusOK, "success", "Brand kit deleted successfully", nil)
}

func Activate(ctx *gin.Context) {
	kitid, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Log("Error while parsing brandkit: " + err.Error())
		setResponse(ctx, http.StatusNotAcceptable, "error", "Invalid brand kit ID", nil)
		return
	}
	if err := brandkit_service.Activate(kitid); err != nil {
		utils.Log("Error while activating brandkit: " + err.Error())
		setResponse(ctx, http.StatusNotAcceptable, "error", "Failed to activate brand kit", nil)
		return
	}

	setResponse(ctx, http.StatusOK, "success", "Brand kit activated successfully", nil)
}

func Duplicate(ctx *gin.Context) {
	kitid, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Log("Error while parsing brandkit: " + err.Error())
		setResponse(ctx, http.StatusNotAcceptable, "error", "Invalid brand kit ID", nil)
		return
	}

	if err := brandkit_service.Duplicate(kitid); err != nil {
		utils.Log("Error while duplicate brandkit: " + err.Error())
		setResponse(ctx, http.StatusNotAcceptable, "error", "Failed to duplicate brand kit", nil)
		return
	}

	setResponse(ctx, http.StatusOK, "success", "Brand kit duplicated successfully", nil)
}

func Search(ctx *gin.Context) {
	var searchReq brandkit_model.SearchRequest
	err := ctx.ShouldBindQuery(&searchReq)

	if err != nil {
		utils.Log("Error while parsing brandkit: " + err.Error())
		utils.SetResponse(ctx, utils.CommonResponse{
			StatusCode: http.StatusFailedDependency,
			Status:     "error",
			Message:    "Failed to search brand kit",
			Data:       nil,
		})
		return
	}

	// Extract fields from searchReq and pass to Search
	brandkits, err := brandkit_service.Search(searchReq.SearchTerm, 0, 10)

	if err != nil {
		utils.Log("Error while searching brandkit: " + err.Error())
		utils.SetResponse(ctx, utils.CommonResponse{
			StatusCode: http.StatusFailedDependency,
			Status:     "error",
			Message:    "Failed to search brand kit",
			Data:       nil,
		})
		return
	}

	// Convert brandkits data into response format, if necessary
	type SearchResponse struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
	}

	var searchResp []SearchResponse
	for _, data := range brandkits.Data {
		searchResp = append(searchResp, SearchResponse{
			Name: data.Name,
			Id:   data.Id,
		})
	}

	utils.SetResponse(ctx, utils.CommonResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "Brand kits searched successfully",
		Data:       searchResp,
	})
}

func Count(ctx *gin.Context) {

	allBrandkits, err := brandkit_service.Count()
	if err != nil {
		utils.Log("Error while counting brandkit: " + err.Error())
		setResponse(ctx, http.StatusAlreadyReported, "error", "Failed to get brand kit count", nil)
		return
	}

	setResponse(ctx, http.StatusOK, "success", "Brand kit count fetched successfully", allBrandkits)
}
