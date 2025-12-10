package divisionController

import (
	"net/http"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/divisionService"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DivisionController struct {
	service divisionService.DivisionService
}

func NewDivisionController(service divisionService.DivisionService) *DivisionController {
	return &DivisionController{service: service}
}

func (c *DivisionController) Create(ctx *gin.Context) {
	var req model.DivisionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}
	
	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	division, err := c.service.CreateDivision(&req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to create division", err, nil)
		return
	}

	apiresponse.Created(ctx, division, "Division created successfully", nil)
}

func (c *DivisionController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	division, err := c.service.GetDivisionByID(id)
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "Division not found", err, nil)
		return
	}

	apiresponse.OK(ctx, division, "Division retrieved successfully", nil)
}

func (c *DivisionController) GetAll(ctx *gin.Context) {
	var filter model.DivisionListRequest
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid query parameters", err, nil)
		return
	}

	if err := validator.New().Struct(&filter); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	divisions, total, err := c.service.GetAllDivisions(&filter)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve divisions", err, nil)
		return
	}

	meta := apiresponse.PageMeta{
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: total,
	}
	apiresponse.OK(ctx, divisions, "Divisions retrieved successfully", meta)
}

func (c *DivisionController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	var req model.DivisionUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	division, err := c.service.UpdateDivision(id, &req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to update division", err, nil)
		return
	}

	apiresponse.OK(ctx, division, "Division updated successfully", nil)
}

func (c *DivisionController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	err = c.service.DeleteDivision(id)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to delete division", err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Division deleted successfully", nil)
}