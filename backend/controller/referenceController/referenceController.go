package referenceController

import (
	"net/http"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/referenceService"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ReferenceController struct {
	referenceService referenceService.ReferenceService
}

func NewReferenceController(referenceService referenceService.ReferenceService) *ReferenceController {
	return &ReferenceController{
		referenceService: referenceService,
	}
}

// CreateReference godoc
// @Summary Create a new reference
// @Tags References
// @Accept json
// @Produce json
// @Param reference body model.ReferenceRequest true "Reference data"
// @Success 201 {object} apiresponse.Response{data=model.ReferenceResponse}
// @Failure 400 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /references [post]
func (c *ReferenceController) CreateReference(ctx *gin.Context) {
	var req model.ReferenceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	result, err := c.referenceService.CreateReference(req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to create reference", err, nil)
		return
	}

	apiresponse.Created(ctx, result, "Reference created successfully", nil)
}

// GetReferences godoc
// @Summary Get all references
// @Tags References
// @Produce json
// @Param reference_category_id query int false "Filter by category ID"
// @Param module_id query int false "Filter by module ID"
// @Param name query string false "Filter by name"
// @Param is_active query boolean false "Filter by active status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} apiresponse.Response{data=[]model.ReferenceResponse}
// @Failure 500 {object} apiresponse.Response
// @Router /references [get]
func (c *ReferenceController) GetReferences(ctx *gin.Context) {
	var req model.ReferenceListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid query parameters", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	result, total, err := c.referenceService.GetReferences(req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve references", err, nil)
		return
	}

	meta := apiresponse.PageMeta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	apiresponse.OK(ctx, result, "References retrieved successfully", meta)
}

// GetReferenceByID godoc
// @Summary Get reference by ID
// @Tags References
// @Produce json
// @Param id path int true "Reference ID"
// @Success 200 {object} apiresponse.Response{data=model.ReferenceResponse}
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /references/{id} [get]
func (c *ReferenceController) GetReferenceByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid reference ID", err, nil)
		return
	}

	result, err := c.referenceService.GetReferenceByID(id)
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "Reference not found", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Reference retrieved successfully", nil)
}

// GetReferencesByCategoryID godoc
// @Summary Get references by category ID
// @Tags References
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} apiresponse.Response{data=[]model.ReferenceResponse}
// @Failure 500 {object} apiresponse.Response
// @Router /reference-categories/{id}/references [get]
func (c *ReferenceController) GetReferencesByCategoryID(ctx *gin.Context) {
	categoryID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid category ID", err, nil)
		return
	}

	result, err := c.referenceService.GetReferencesByCategoryID(categoryID)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve references", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "References retrieved successfully", nil)
}

// GetReferencesByModuleID godoc
// @Summary Get references by module ID
// @Tags References
// @Produce json
// @Param id path int true "Module ID"
// @Success 200 {object} apiresponse.Response{data=[]model.ReferenceResponse}
// @Failure 500 {object} apiresponse.Response
// @Router /modules/{id}/references [get]
func (c *ReferenceController) GetReferencesByModuleID(ctx *gin.Context) {
	moduleID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid module ID", err, nil)
		return
	}

	result, err := c.referenceService.GetReferencesByModuleID(moduleID)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve references", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "References retrieved successfully", nil)
}

// UpdateReference godoc
// @Summary Update a reference
// @Tags References
// @Accept json
// @Produce json
// @Param id path int true "Reference ID"
// @Param reference body model.ReferenceRequest true "Reference data"
// @Success 200 {object} apiresponse.Response{data=model.ReferenceResponse}
// @Failure 400 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /references/{id} [put]
func (c *ReferenceController) UpdateReference(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid reference ID", err, nil)
		return
	}

	var req model.ReferenceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	result, err := c.referenceService.UpdateReference(id, req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to update reference", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Reference updated successfully", nil)
}

// DeleteReference godoc
// @Summary Delete a reference
// @Tags References
// @Produce json
// @Param id path int true "Reference ID"
// @Success 200 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /references/{id} [delete]
func (c *ReferenceController) DeleteReference(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid reference ID", err, nil)
		return
	}

	if err := c.referenceService.DeleteReference(id); err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to delete reference", err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Reference deleted successfully", nil)
}
