package roleController

import (
	"net/http"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/roleService"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RoleController struct {
	service roleService.RoleService
}

func NewRoleController(service roleService.RoleService) *RoleController {
	return &RoleController{service: service}
}

func (c *RoleController) Create(ctx *gin.Context) {
	var req model.RoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	role, err := c.service.CreateRole(&req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to create role", err, nil)
		return
	}

	apiresponse.Created(ctx, role, "Role created successfully", nil)
}

func (c *RoleController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	role, err := c.service.GetRoleByID(id)
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "Role not found", err, nil)
		return
	}

	apiresponse.OK(ctx, role, "Role retrieved successfully", nil)
}

func (c *RoleController) GetAll(ctx *gin.Context) {
	var filter model.RoleListRequest
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid query parameters", err, nil)
		return
	}

	if err := validator.New().Struct(&filter); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	roles, total, err := c.service.GetAllRoles(&filter)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve roles", err, nil)
		return
	}

	meta := apiresponse.PageMeta{
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: total,
	}
	apiresponse.OK(ctx, roles, "Roles retrieved successfully", meta)
}

func (c *RoleController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	var req model.RoleUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	role, err := c.service.UpdateRole(id, &req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to update role", err, nil)
		return
	}

	apiresponse.OK(ctx, role, "Role updated successfully", nil)
}

func (c *RoleController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	err = c.service.DeleteRole(id)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to delete role", err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Role deleted successfully", nil)
}
