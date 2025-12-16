package moduleController

import (
	"net/http"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/moduleService"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ModuleController struct {
	moduleService moduleService.ModuleService
}

func NewModuleController(moduleService moduleService.ModuleService) *ModuleController {
	return &ModuleController{
		moduleService: moduleService,
	}
}

// CreateModule godoc
// @Summary Create a new module
// @Tags Modules
// @Accept json
// @Produce json
// @Param module body model.ModuleRequest true "Module data"
// @Success 201 {object} apiresponse.Response{data=model.ModuleResponse}
// @Failure 400 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /modules [post]
func (c *ModuleController) CreateModule(ctx *gin.Context) {
	var req model.ModuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	result, err := c.moduleService.CreateModule(req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to create module", err, nil)
		return
	}

	apiresponse.Created(ctx, result, "Module created successfully", nil)
}

// GetModules godoc
// @Summary Get all modules
// @Tags Modules
// @Produce json
// @Param code query string false "Filter by code"
// @Param name query string false "Filter by name"
// @Param is_active query boolean false "Filter by active status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} apiresponse.Response{data=[]model.ModuleResponse}
// @Failure 500 {object} apiresponse.Response
// @Router /modules [get]
func (c *ModuleController) GetModules(ctx *gin.Context) {
	var req model.ModuleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid query parameters", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	result, total, err := c.moduleService.GetModules(req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve modules", err, nil)
		return
	}

	meta := apiresponse.PageMeta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	apiresponse.OK(ctx, result, "Modules retrieved successfully", meta)
}

// GetModuleByID godoc
// @Summary Get module by ID
// @Tags Modules
// @Produce json
// @Param id path int true "Module ID"
// @Success 200 {object} apiresponse.Response{data=model.ModuleResponse}
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /modules/{id} [get]
func (c *ModuleController) GetModuleByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid module ID", err, nil)
		return
	}

	result, err := c.moduleService.GetModuleByID(id)
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "Module not found", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Module retrieved successfully", nil)
}

// UpdateModule godoc
// @Summary Update a module
// @Tags Modules
// @Accept json
// @Produce json
// @Param id path int true "Module ID"
// @Param module body model.ModuleRequest true "Module data"
// @Success 200 {object} apiresponse.Response{data=model.ModuleResponse}
// @Failure 400 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /modules/{id} [put]
func (c *ModuleController) UpdateModule(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid module ID", err, nil)
		return
	}

	var req model.ModuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	result, err := c.moduleService.UpdateModule(id, req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to update module", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Module updated successfully", nil)
}

// DeleteModule godoc
// @Summary Delete a module
// @Tags Modules
// @Produce json
// @Param id path int true "Module ID"
// @Success 200 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /modules/{id} [delete]
func (c *ModuleController) DeleteModule(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid module ID", err, nil)
		return
	}

	if err := c.moduleService.DeleteModule(id); err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to delete module", err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Module deleted successfully", nil)
}
