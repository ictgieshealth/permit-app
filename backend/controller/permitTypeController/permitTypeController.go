package permitTypeController

import (
	"net/http"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/permitTypeService"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PermitTypeController struct {
	service permitTypeService.PermitTypeService
}

func NewPermitTypeController(service permitTypeService.PermitTypeService) *PermitTypeController {
	return &PermitTypeController{service: service}
}

func (c *PermitTypeController) Create(ctx *gin.Context) {
	var req model.PermitTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	permitType, err := c.service.CreatePermitType(&req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to create permit type", err, nil)
		return
	}

	apiresponse.Created(ctx, permitType, "Permit type created successfully", nil)
}

func (c *PermitTypeController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	permitType, err := c.service.GetPermitTypeByID(id)
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "Permit type not found", err, nil)
		return
	}

	apiresponse.OK(ctx, permitType, "Permit type retrieved successfully", nil)
}

func (c *PermitTypeController) GetAll(ctx *gin.Context) {
	var filter model.PermitTypeListRequest
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid query parameters", err, nil)
		return
	}

	if err := validator.New().Struct(&filter); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	permitTypes, total, err := c.service.GetAllPermitTypes(&filter)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve permit types", err, nil)
		return
	}

	meta := apiresponse.PageMeta{
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: total,
	}
	apiresponse.OK(ctx, permitTypes, "Permit types retrieved successfully", meta)
}

func (c *PermitTypeController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	var req model.PermitTypeUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	permitType, err := c.service.UpdatePermitType(id, &req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to update permit type", err, nil)
		return
	}

	apiresponse.OK(ctx, permitType, "Permit type updated successfully", nil)
}

func (c *PermitTypeController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	err = c.service.DeletePermitType(id)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to delete permit type", err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Permit type deleted successfully", nil)
}