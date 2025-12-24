package taskController

import (
	"net/http"
	"permit-app/helper/apiRequest"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/taskService"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TaskController struct {
	taskService taskService.TaskService
}

func NewTaskController(taskService taskService.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

// Create creates a new task
func (c *TaskController) Create(ctx *gin.Context) {
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists || userID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "User context not found", nil, nil)
		return
	}

	// Parse form data
	projectID, err := strconv.ParseInt(ctx.PostForm("project_id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid project_id", err, nil)
		return
	}

	title := ctx.PostForm("title")
	if title == "" {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Title is required", nil, nil)
		return
	}

	priorityID, err := strconv.ParseInt(ctx.PostForm("priority_id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid priority_id", err, nil)
		return
	}

	req := &model.TaskRequest{
		ProjectID:  projectID,
		Title:      title,
		PriorityID: priorityID,
	}

	description := ctx.PostForm("description")
	if description != "" {
		req.Description = &description
	}

	if assignedIDStr := ctx.PostForm("assigned_id"); assignedIDStr != "" {
		if assignedID, err := strconv.ParseInt(assignedIDStr, 10, 64); err == nil {
			req.AssignedID = &assignedID
		}
	}

	if stackIDStr := ctx.PostForm("stack_id"); stackIDStr != "" {
		if stackID, err := strconv.ParseInt(stackIDStr, 10, 64); err == nil {
			req.StackID = &stackID
		}
	}

	dueDate := ctx.PostForm("due_date")
	if dueDate != "" {
		req.DueDate = &dueDate
	}

	if err := validator.New().Struct(req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	// Handle file uploads
	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	task, err := c.taskService.Create(req, files, domainID.(int64), userID.(int64))
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to create task", err, nil)
		return
	}

	apiresponse.Created(ctx, task, "Task created successfully", nil)
}

// GetAll retrieves all tasks with filters
func (c *TaskController) GetAll(ctx *gin.Context) {
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}

	req := &model.TaskListRequest{
		Page:  apiRequest.ParseInt(ctx, "page", 1),
		Limit: apiRequest.ParseInt(ctx, "limit", 10),
	}

	if search := ctx.Query("search"); search != "" {
		req.Search = &search
	}

	if projectIDStr := ctx.Query("project_id"); projectIDStr != "" {
		if projectID, err := strconv.ParseInt(projectIDStr, 10, 64); err == nil {
			req.ProjectID = &projectID
		}
	}

	if statusIDStr := ctx.Query("status_id"); statusIDStr != "" {
		if statusID, err := strconv.ParseInt(statusIDStr, 10, 64); err == nil {
			req.StatusID = &statusID
		}
	}

	if approvalStatusIDStr := ctx.Query("approval_status_id"); approvalStatusIDStr != "" {
		if approvalStatusID, err := strconv.ParseInt(approvalStatusIDStr, 10, 64); err == nil {
			req.ApprovalStatusID = &approvalStatusID
		}
	}

	if assignedIDStr := ctx.Query("assigned_id"); assignedIDStr != "" {
		if assignedID, err := strconv.ParseInt(assignedIDStr, 10, 64); err == nil {
			req.AssignedID = &assignedID
		}
	}

	if startDate := ctx.Query("start_date"); startDate != "" {
		req.StartDate = &startDate
	}

	if endDate := ctx.Query("end_date"); endDate != "" {
		req.EndDate = &endDate
	}

	tasks, total, err := c.taskService.GetAll(domainID.(int64), req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve tasks", err, nil)
		return
	}

	meta := apiresponse.PageMeta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	apiresponse.OK(ctx, tasks, "Tasks retrieved successfully", meta)
}

// GetByID retrieves a task by ID
func (c *TaskController) GetByID(ctx *gin.Context) {
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid task ID", err, nil)
		return
	}

	task, err := c.taskService.GetByID(id, domainID.(int64))
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "Task not found", err, nil)
		return
	}

	apiresponse.OK(ctx, task, "Task retrieved successfully", nil)
}

// GetByCode retrieves a task by code
func (c *TaskController) GetByCode(ctx *gin.Context) {
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}

	code := ctx.Param("code")
	if code == "" {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Task code is required", nil, nil)
		return
	}

	task, err := c.taskService.GetByCode(code, domainID.(int64))
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "Task not found", err, nil)
		return
	}

	apiresponse.OK(ctx, task, "Task retrieved successfully", nil)
}

// Update updates a task
func (c *TaskController) Update(ctx *gin.Context) {
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists || userID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "User context not found", nil, nil)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid task ID", err, nil)
		return
	}

	title := ctx.PostForm("title")
	if title == "" {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Title is required", nil, nil)
		return
	}

	projectID, err := strconv.ParseInt(ctx.PostForm("project_id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid project_id", err, nil)
		return
	}

	priorityID, err := strconv.ParseInt(ctx.PostForm("priority_id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid priority_id", err, nil)
		return
	}

	req := &model.TaskUpdateRequest{
		Title:      title,
		ProjectID:  projectID,
		PriorityID: priorityID,
	}

	if description := ctx.PostForm("description"); description != "" {
		req.Description = &description
	}

	if descBefore := ctx.PostForm("description_before"); descBefore != "" {
		req.DescriptionBefore = &descBefore
	}

	if descAfter := ctx.PostForm("description_after"); descAfter != "" {
		req.DescriptionAfter = &descAfter
	}

	if assignedIDStr := ctx.PostForm("assigned_id"); assignedIDStr != "" {
		if assignedID, err := strconv.ParseInt(assignedIDStr, 10, 64); err == nil {
			req.AssignedID = &assignedID
		}
	}

	if stackIDStr := ctx.PostForm("stack_id"); stackIDStr != "" {
		if stackID, err := strconv.ParseInt(stackIDStr, 10, 64); err == nil {
			req.StackID = &stackID
		}
	}

	if dueDate := ctx.PostForm("due_date"); dueDate != "" {
		req.DueDate = &dueDate
	}

	if err := validator.New().Struct(req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	form, _ := ctx.MultipartForm()
	files := form.File["files"]
	
	// Get deleted file IDs
	var deletedFileIds []int64
	if deletedIdsStr := ctx.PostFormArray("deleted_file_ids[]"); len(deletedIdsStr) > 0 {
		for _, idStr := range deletedIdsStr {
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				deletedFileIds = append(deletedFileIds, id)
			}
		}
	}

	task, err := c.taskService.Update(id, req, files, deletedFileIds, domainID.(int64), userID.(int64))
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to update task", err, nil)
		return
	}

	apiresponse.OK(ctx, task, "Task updated successfully", nil)
}

// Delete soft deletes a task
func (c *TaskController) Delete(ctx *gin.Context) {
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid task ID", err, nil)
		return
	}

	err = c.taskService.Delete(id, domainID.(int64))
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to delete task", err, nil)
		return
	}

	apiresponse.OK(ctx, map[string]interface{}{}, "Task deleted successfully", nil)
}

// ChangeStatus changes task status
func (c *TaskController) ChangeStatus(ctx *gin.Context) {
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists || userID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "User context not found", nil, nil)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid task ID", err, nil)
		return
	}

	var req model.TaskChangeStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	err = c.taskService.ChangeStatus(id, &req, domainID.(int64), userID.(int64))
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to change task status", err, nil)
		return
	}

	apiresponse.OK(ctx, map[string]interface{}{}, "Task status changed successfully", nil)
}

// ChangeType changes task type
func (c *TaskController) ChangeType(ctx *gin.Context) {
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists || userID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "User context not found", nil, nil)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid task ID", err, nil)
		return
	}

	var req model.TaskChangeTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	err = c.taskService.ChangeType(id, &req, domainID.(int64), userID.(int64))
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to change task type", err, nil)
		return
	}

	apiresponse.OK(ctx, map[string]interface{}{}, "Task type changed successfully", nil)
}

// InReview sets task to in review status
func (c *TaskController) InReview(ctx *gin.Context) {
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists || userID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "User context not found", nil, nil)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid task ID", err, nil)
		return
	}

	req := &model.TaskInReviewRequest{}

	if descBefore := ctx.PostForm("description_before"); descBefore != "" {
		req.DescriptionBefore = &descBefore
	}

	if descAfter := ctx.PostForm("description_after"); descAfter != "" {
		req.DescriptionAfter = &descAfter
	}

	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	err = c.taskService.InReview(id, req, files, domainID.(int64), userID.(int64))
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to set task in review", err, nil)
		return
	}

	apiresponse.OK(ctx, map[string]interface{}{}, "Task set to in review successfully", nil)
}

// SetReason sets task reason
func (c *TaskController) SetReason(ctx *gin.Context) {
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists || userID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "User context not found", nil, nil)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid task ID", err, nil)
		return
	}

	var req model.TaskReasonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	err = c.taskService.SetReason(id, &req, domainID.(int64), userID.(int64))
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to set task reason", err, nil)
		return
	}

	apiresponse.OK(ctx, map[string]interface{}{}, "Task reason set successfully", nil)
}

// SetRevision sets task revision
func (c *TaskController) SetRevision(ctx *gin.Context) {
	domainID, exists := ctx.Get("domain_id")
	if !exists || domainID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Domain context not found", nil, nil)
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists || userID == nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "User context not found", nil, nil)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid task ID", err, nil)
		return
	}

	req := &model.TaskRevisionRequest{}

	if revision := ctx.PostForm("revision"); revision != "" {
		req.Revision = &revision
	}

	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	err = c.taskService.SetRevision(id, req, files, domainID.(int64), userID.(int64))
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to set task revision", err, nil)
		return
	}

	apiresponse.OK(ctx, map[string]interface{}{}, "Task revision set successfully", nil)
}
