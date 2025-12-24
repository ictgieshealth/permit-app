package projectController

import (
	"net/http"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/projectService"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProjectController struct {
	projectService projectService.ProjectService
}

func NewProjectController(projectService projectService.ProjectService) *ProjectController {
	return &ProjectController{
		projectService: projectService,
	}
}

// CreateProject godoc
// @Summary Create a new project
// @Tags Projects
// @Accept json
// @Produce json
// @Param project body model.ProjectRequest true "Project data"
// @Success 201 {object} apiresponse.Response{data=model.ProjectResponse}
// @Failure 400 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /projects [post]
func (c *ProjectController) CreateProject(ctx *gin.Context) {
	var req model.ProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	// Validate domain_id from JWT token
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}

	if req.DomainID != domainID.(int64) {
		apiresponse.Error(ctx, http.StatusForbidden, "FORBIDDEN", "Cannot create project in different domain", nil, nil)
		return
	}

	result, err := c.projectService.CreateProject(req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to create project", err, nil)
		return
	}

	apiresponse.Created(ctx, result, "Project created successfully", nil)
}

// GetProjects godoc
// @Summary Get all projects
// @Tags Projects
// @Produce json
// @Param domain_id query int false "Filter by domain ID"
// @Param project_status_id query int false "Filter by project status ID"
// @Param name query string false "Filter by name"
// @Param code query string false "Filter by code"
// @Param status query boolean false "Filter by status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} apiresponse.Response{data=[]model.ProjectResponse}
// @Failure 500 {object} apiresponse.Response
// @Router /projects [get]
func (c *ProjectController) GetProjects(ctx *gin.Context) {
	var req model.ProjectListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid query parameters", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	// Extract domain_id from JWT token context
	if domainID, exists := ctx.Get("domain_id"); exists && domainID != nil {
		req.DomainID = domainID.(int64)
	}

	result, total, err := c.projectService.GetProjects(req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve projects", err, nil)
		return
	}

	meta := apiresponse.PageMeta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	apiresponse.OK(ctx, result, "Projects retrieved successfully", meta)
}

// GetProjectByID godoc
// @Summary Get project by ID
// @Tags Projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} apiresponse.Response{data=model.ProjectResponse}
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /projects/{id} [get]
func (c *ProjectController) GetProjectByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid project ID", err, nil)
		return
	}

	result, err := c.projectService.GetProjectByID(id)
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "Project not found", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Project retrieved successfully", nil)
}

// GetProjectsByDomainID godoc
// @Summary Get projects by domain ID
// @Tags Projects
// @Produce json
// @Param id path int true "Domain ID"
// @Success 200 {object} apiresponse.Response{data=[]model.ProjectResponse}
// @Failure 500 {object} apiresponse.Response
// @Router /domains/{id}/projects [get]
func (c *ProjectController) GetProjectsByDomainID(ctx *gin.Context) {
	domainID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid domain ID", err, nil)
		return
	}

	result, err := c.projectService.GetProjectsByDomainID(domainID)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve projects", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Projects retrieved successfully", nil)
}

// GetProjectsByUserID godoc
// @Summary Get projects by user ID
// @Tags Projects
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} apiresponse.Response{data=[]model.ProjectResponse}
// @Failure 500 {object} apiresponse.Response
// @Router /users/{id}/projects [get]
func (c *ProjectController) GetProjectsByUserID(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid user ID", err, nil)
		return
	}

	result, err := c.projectService.GetProjectsByUserID(userID)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve projects", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Projects retrieved successfully", nil)
}

func (c *ProjectController) GetUsersByProjectID(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid user ID", err, nil)
		return
	}

	result, err := c.projectService.GetUsersByProjectID(userID)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve projects", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Projects retrieved successfully", nil)
}

// UpdateProject godoc
// @Summary Update a project
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body model.ProjectUpdateRequest true "Project data"
// @Success 200 {object} apiresponse.Response{data=model.ProjectResponse}
// @Failure 400 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /projects/{id} [put]
func (c *ProjectController) UpdateProject(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid project ID", err, nil)
		return
	}

	var req model.ProjectUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	result, err := c.projectService.UpdateProject(id, req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to update project", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Project updated successfully", nil)
}

// DeleteProject godoc
// @Summary Delete a project
// @Tags Projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /projects/{id} [delete]
func (c *ProjectController) DeleteProject(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid project ID", err, nil)
		return
	}

	if err := c.projectService.DeleteProject(id); err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, err.Error(), err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Project deleted successfully", nil)
}

// ChangeProjectStatus godoc
// @Summary Change project status
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param status body model.ProjectStatusChangeRequest true "Status data"
// @Success 200 {object} apiresponse.Response{data=model.ProjectResponse}
// @Failure 400 {object} apiresponse.Response
// @Failure 404 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Router /projects/{id}/change-status [post]
func (c *ProjectController) ChangeProjectStatus(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid project ID", err, nil)
		return
	}

	var req model.ProjectStatusChangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	result, err := c.projectService.ChangeProjectStatus(id, req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to change project status", err, nil)
		return
	}

	apiresponse.OK(ctx, result, "Project status updated successfully", nil)
}
