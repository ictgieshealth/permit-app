package referenceCategoryController

import (
	"net/http"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/referenceCategoryService"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ReferenceCategoryController struct {
	categoryService referenceCategoryService.ReferenceCategoryService
}

func NewReferenceCategoryController(categoryService referenceCategoryService.ReferenceCategoryService) *ReferenceCategoryController {
	return &ReferenceCategoryController{
		categoryService: categoryService,
	}
}

// CreateReferenceCategory godoc
// @Summary Create a new reference category
// @Tags Reference Categories
// @Accept json
// @Produce json
// @Param category body model.ReferenceCategoryRequest true "Category data"
// @Success 201 {object} apiresponse.Response{data=model.ReferenceCategoryResponse}
// @Failure 400 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /reference-categories [post]
func (c *ReferenceCategoryController) CreateReferenceCategory(ctx *gin.Context) {
	var req model.ReferenceCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	result, err := c.categoryService.CreateReferenceCategory(req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to create category", err, nil)
		return
	}

	apiresponse.Created(ctx, result, "Reference category created successfully", nil)
}

// GetReferenceCategories godoc
// @Summary Get all reference categories
// @Tags Reference Categories
// @Produce json
// @Param module_id query int false "Filter by module ID"
// @Param name query string false "Filter by name"
// @Param is_active query boolean false "Filter by active status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} apiresponse.Response{data=[]model.ReferenceCategoryResponse}
// @Failure 500 {object} apiresponse.Response
// @Router /reference-categories [get]
func (c *ReferenceCategoryController) GetReferenceCategories(ctx *gin.Context) {
	var req model.ReferenceCategoryListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid query parameters", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	result, total, err := c.categoryService.GetReferenceCategories(req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve categories", err, nil)
		return
	}

	meta := apiresponse.PageMeta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	apiresponse.OK(ctx, result, "Reference categories retrieved successfully", meta)
}

// GetReferenceCategoryByID godoc
// @Summary Get reference category by ID
// @Tags Reference Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} apiresponse.Response{data=model.ReferenceCategoryResponse}
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /reference-categories/{id} [get]
func (c *ReferenceCategoryController) GetReferenceCategoryByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid category ID", err, nil)
		return
	}

	result, err := c.categoryService.GetReferenceCategoryByID(id)
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "Category not found", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Reference category retrieved successfully", nil)
}

// GetCategoriesByModuleID godoc
// @Summary Get categories by module ID
// @Tags Reference Categories
// @Produce json
// @Param id path int true "Module ID"
// @Success 200 {object} apiresponse.Response{data=[]model.ReferenceCategoryResponse}
// @Failure 500 {object} apiresponse.Response
// @Router /modules/{id}/categories [get]
func (c *ReferenceCategoryController) GetCategoriesByModuleID(ctx *gin.Context) {
	moduleID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid module ID", err, nil)
		return
	}

	result, err := c.categoryService.GetCategoriesByModuleID(moduleID)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve categories", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Categories retrieved successfully", nil)
}

// UpdateReferenceCategory godoc
// @Summary Update a reference category
// @Tags Reference Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body model.ReferenceCategoryRequest true "Category data"
// @Success 200 {object} apiresponse.Response{data=model.ReferenceCategoryResponse}
// @Failure 400 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /reference-categories/{id} [put]
func (c *ReferenceCategoryController) UpdateReferenceCategory(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid category ID", err, nil)
		return
	}

	var req model.ReferenceCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	result, err := c.categoryService.UpdateReferenceCategory(id, req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to update category", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Reference category updated successfully", nil)
}

// DeleteReferenceCategory godoc
// @Summary Delete a reference category
// @Tags Reference Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /reference-categories/{id} [delete]
func (c *ReferenceCategoryController) DeleteReferenceCategory(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid category ID", err, nil)
		return
	}

	if err := c.categoryService.DeleteReferenceCategory(id); err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to delete category", err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Reference category deleted successfully", nil)
}
