package taskRequestController

import (
	"net/http"
	"permit-app/helper/apiRequest"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/taskService"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskRequestController struct {
	taskService taskService.TaskService
}

func NewTaskRequestController(taskService taskService.TaskService) *TaskRequestController {
	return &TaskRequestController{
		taskService: taskService,
	}
}

// GetAll retrieves all task approval requests
// This matches PHP TaskRequestController::index() logic
func (c *TaskRequestController) GetAll(ctx *gin.Context) {
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

	// Use GetAllRequests which doesn't force approval_status_id filter
	tasks, total, err := c.taskService.GetAllRequests(domainID.(int64), req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve task requests", err, nil)
		return
	}

	meta := apiresponse.PageMeta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	apiresponse.OK(ctx, tasks, "Task requests retrieved successfully", meta)
}

// ApproveTask approves a task in the approval workflow
func (c *TaskRequestController) ApproveTask(ctx *gin.Context) {
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

	taskID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid task ID", err, nil)
		return
	}

	approvalTaskID, err := strconv.ParseInt(ctx.Param("approval_id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid approval task ID", err, nil)
		return
	}

	var req model.ApprovalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Note is optional, so we ignore binding errors
		req.Note = nil
	}

	err = c.taskService.ApproveTask(taskID, approvalTaskID, &req, domainID.(int64), userID.(int64))
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to approve task", err, nil)
		return
	}

	apiresponse.OK(ctx, map[string]interface{}{}, "Task approved successfully", nil)
}

// RejectTask rejects a task in the approval workflow
func (c *TaskRequestController) RejectTask(ctx *gin.Context) {
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

	taskID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid task ID", err, nil)
		return
	}

	approvalTaskID, err := strconv.ParseInt(ctx.Param("approval_id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid approval task ID", err, nil)
		return
	}

	var req model.ApprovalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Note is optional, so we ignore binding errors
		req.Note = nil
	}

	err = c.taskService.RejectTask(taskID, approvalTaskID, &req, domainID.(int64), userID.(int64))
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to reject task", err, nil)
		return
	}

	apiresponse.OK(ctx, map[string]interface{}{}, "Task rejected successfully", nil)
}
