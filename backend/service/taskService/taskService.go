package taskService

import (
	"errors"
	"fmt"
	"mime/multipart"
	"permit-app/helper"
	"permit-app/model"
	"permit-app/repo/taskRepository"
	"time"
)

type TaskService interface {
	Create(req *model.TaskRequest, files []*multipart.FileHeader, domainID, userID int64) (*model.Task, error)
	GetByID(id int64, domainID int64) (*model.TaskResponse, error)
	GetByCode(code string, domainID int64) (*model.TaskResponse, error)
	GetAll(domainID int64, filters *model.TaskListRequest) ([]model.TaskResponse, int64, error)
	GetAllRequests(domainID int64, filters *model.TaskListRequest) ([]model.TaskResponse, int64, error)
	Update(id int64, req *model.TaskUpdateRequest, files []*multipart.FileHeader, domainID, userID int64) (*model.Task, error)
	Delete(id int64, domainID int64) error
	ChangeStatus(id int64, req *model.TaskChangeStatusRequest, domainID, userID int64) error
	ChangeType(id int64, req *model.TaskChangeTypeRequest, domainID, userID int64) error
	InReview(id int64, req *model.TaskInReviewRequest, files []*multipart.FileHeader, domainID, userID int64) error
	SetReason(id int64, req *model.TaskReasonRequest, domainID, userID int64) error
	SetRevision(id int64, req *model.TaskRevisionRequest, files []*multipart.FileHeader, domainID, userID int64) error

	// Approval
	ApproveTask(taskID, approvalTaskID int64, req *model.ApprovalRequest, domainID, userID int64) error
	RejectTask(taskID, approvalTaskID int64, req *model.ApprovalRequest, domainID, userID int64) error
}

type taskService struct {
	taskRepo taskRepository.TaskRepository
}

func NewTaskService(taskRepo taskRepository.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) Create(req *model.TaskRequest, files []*multipart.FileHeader, domainID, userID int64) (*model.Task, error) {
	// Generate task code
	code, err := s.taskRepo.GenerateCode(req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate task code: %v", err)
	}

	// Parse due date
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err == nil {
			dueDate = &parsed
		}
	}

	// Set default status
	statusID := int64(helper.TaskStatusToDo)
	approvalStatusID := int64(helper.ApprovalStatusWaiting)

	task := &model.Task{
		DomainID:         domainID,
		ProjectID:        req.ProjectID,
		Code:             code,
		Title:            req.Title,
		Description:      req.Description,
		StatusID:         &statusID,
		PriorityID:       &req.PriorityID,
		AssignedID:       req.AssignedID,
		StackID:          req.StackID,
		ApprovalStatusID: &approvalStatusID,
		DueDate:          dueDate,
		CreatedBy:        userID,
		Status:           true,
	}

	// Create task
	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %v", err)
	}

	// Create approval tasks (2 sequences)
	approvalTasks := []model.ApprovalTask{
		{
			TaskID:           task.ID,
			Sequence:         1,
			ApprovalStatusID: ptrInt64(helper.ApprovalStatusWaiting),
			Status:           true,
		},
		{
			TaskID:           task.ID,
			Sequence:         2,
			ApprovalStatusID: ptrInt64(helper.ApprovalStatusWaiting),
			Status:           true,
		},
	}

	if err := s.taskRepo.CreateApprovalTasks(approvalTasks); err != nil {
		return nil, fmt.Errorf("failed to create approval tasks: %v", err)
	}

	// Handle file uploads
	if len(files) > 0 {
		if err := s.uploadTaskFiles(task.ID, files, 30); err != nil { // 30 = Create file type
			return nil, fmt.Errorf("failed to upload files: %v", err)
		}
	}

	return task, nil
}

func (s *taskService) GetByID(id int64, domainID int64) (*model.TaskResponse, error) {
	task, err := s.taskRepo.GetByID(id, domainID)
	if err != nil {
		return nil, err
	}

	return s.toTaskResponse(task), nil
}

func (s *taskService) GetByCode(code string, domainID int64) (*model.TaskResponse, error) {
	task, err := s.taskRepo.GetByCode(code, domainID)
	if err != nil {
		return nil, err
	}

	return s.toTaskResponse(task), nil
}

func (s *taskService) GetAll(domainID int64, filters *model.TaskListRequest) ([]model.TaskResponse, int64, error) {
	filterMap := make(map[string]interface{})

	if filters.Search != nil {
		filterMap["search"] = *filters.Search
	}
	if filters.ProjectID != nil {
		filterMap["project_id"] = *filters.ProjectID
	}
	if filters.StatusID != nil {
		filterMap["status_id"] = *filters.StatusID
	}
	// Force approval_status_id = 22 (Approved) for task list, matching PHP TaskController logic
	filterMap["approval_status_id"] = int64(helper.ApprovalStatusApprove)

	if filters.AssignedID != nil {
		filterMap["assigned_id"] = *filters.AssignedID
	}
	if filters.StartDate != nil {
		filterMap["start_date"] = *filters.StartDate
	}
	if filters.EndDate != nil {
		filterMap["end_date"] = *filters.EndDate
	}

	tasks, total, err := s.taskRepo.GetAll(domainID, filterMap, filters.Page, filters.Limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]model.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = *s.toTaskResponse(&task)
	}

	return responses, total, nil
}

// GetAllRequests retrieves all task approval requests without filtering by approval status
// This matches PHP TaskRequestController logic
func (s *taskService) GetAllRequests(domainID int64, filters *model.TaskListRequest) ([]model.TaskResponse, int64, error) {
	filterMap := make(map[string]interface{})

	if filters.Search != nil {
		filterMap["search"] = *filters.Search
	}
	if filters.ProjectID != nil {
		filterMap["project_id"] = *filters.ProjectID
	}
	if filters.StatusID != nil {
		filterMap["status_id"] = *filters.StatusID
	}
	if filters.ApprovalStatusID != nil {
		filterMap["approval_status_id"] = *filters.ApprovalStatusID
	}
	if filters.AssignedID != nil {
		filterMap["assigned_id"] = *filters.AssignedID
	}
	if filters.StartDate != nil {
		filterMap["start_date"] = *filters.StartDate
	}
	if filters.EndDate != nil {
		filterMap["end_date"] = *filters.EndDate
	}

	tasks, total, err := s.taskRepo.GetAll(domainID, filterMap, filters.Page, filters.Limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]model.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = *s.toTaskResponse(&task)
	}

	return responses, total, nil
}

func (s *taskService) Update(id int64, req *model.TaskUpdateRequest, files []*multipart.FileHeader, domainID, userID int64) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(id, domainID)
	if err != nil {
		return nil, err
	}

	// Parse due date
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err == nil {
			dueDate = &parsed
		}
	}

	// Update fields
	task.Title = req.Title
	task.Description = req.Description
	task.DescriptionBefore = req.DescriptionBefore
	task.DescriptionAfter = req.DescriptionAfter
	task.ProjectID = req.ProjectID
	task.PriorityID = &req.PriorityID
	task.AssignedID = req.AssignedID
	task.StackID = req.StackID
	task.DueDate = dueDate
	task.UpdatedBy = &userID

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	// Handle file uploads if any
	if len(files) > 0 {
		if err := s.uploadTaskFiles(task.ID, files, 30); err != nil {
			return nil, err
		}
	}

	return task, nil
}

func (s *taskService) Delete(id int64, domainID int64) error {
	return s.taskRepo.Delete(id, domainID)
}

func (s *taskService) ChangeStatus(id int64, req *model.TaskChangeStatusRequest, domainID, userID int64) error {
	return s.taskRepo.ChangeStatus(id, domainID, req.StatusID, userID)
}

func (s *taskService) ChangeType(id int64, req *model.TaskChangeTypeRequest, domainID, userID int64) error {
	return s.taskRepo.ChangeType(id, domainID, req.TypeID, userID)
}

func (s *taskService) InReview(id int64, req *model.TaskInReviewRequest, files []*multipart.FileHeader, domainID, userID int64) error {
	if err := s.taskRepo.InReview(id, domainID, req.DescriptionBefore, req.DescriptionAfter, userID); err != nil {
		return err
	}

	// Upload before/after files
	if len(files) > 0 {
		if err := s.uploadTaskFiles(id, files, 31); err != nil { // 31 = Before file type
			return err
		}
	}

	// Change status to In Review
	return s.taskRepo.ChangeStatus(id, domainID, helper.TaskStatusInReview, userID)
}

func (s *taskService) SetReason(id int64, req *model.TaskReasonRequest, domainID, userID int64) error {
	return s.taskRepo.SetReason(id, domainID, req.Reason, userID)
}

func (s *taskService) SetRevision(id int64, req *model.TaskRevisionRequest, files []*multipart.FileHeader, domainID, userID int64) error {
	if err := s.taskRepo.SetRevision(id, domainID, req.Revision, userID); err != nil {
		return err
	}

	// Upload revision files
	if len(files) > 0 {
		if err := s.uploadTaskFiles(id, files, 38); err != nil { // 38 = Revision file type
			return err
		}
	}

	// Change status to Revision
	return s.taskRepo.ChangeStatus(id, domainID, helper.TaskStatusRevision, userID)
}

func (s *taskService) ApproveTask(taskID, approvalTaskID int64, req *model.ApprovalRequest, domainID, userID int64) error {
	// Get task
	task, err := s.taskRepo.GetByID(taskID, domainID)
	if err != nil {
		return err
	}

	// Get approval tasks
	approvalTasks, err := s.taskRepo.GetApprovalTasksByTaskID(taskID)
	if err != nil {
		return err
	}

	if len(approvalTasks) == 0 {
		return errors.New("no approval tasks found")
	}

	// Find the approval task to approve
	var currentApproval *model.ApprovalTask
	for i := range approvalTasks {
		if approvalTasks[i].ID == approvalTaskID {
			currentApproval = &approvalTasks[i]
			break
		}
	}

	if currentApproval == nil {
		return errors.New("approval task not found")
	}

	// Update approval task
	now := time.Now()
	currentApproval.ApprovedBy = &userID
	currentApproval.ApprovalStatusID = ptrInt64(helper.ApprovalStatusApprove)
	currentApproval.ApprovalDate = &now
	currentApproval.Note = req.Note

	if err := s.taskRepo.UpdateApprovalTask(currentApproval); err != nil {
		return err
	}

	// Check which sequence was approved
	if currentApproval.Sequence == 1 {
		// Sequence 1 approved - set status to Pending Manager
		approvalStatusID := int64(helper.ApprovalStatusPendingManager)
		task.ApprovalStatusID = &approvalStatusID
	} else if currentApproval.Sequence == 2 {
		// Sequence 2 approved - task fully approved
		approvalStatusID := int64(helper.ApprovalStatusApprove)
		task.ApprovalStatusID = &approvalStatusID
		task.ApprovedBy = &userID
		task.ApprovalDate = &now
	}

	task.UpdatedBy = &userID
	return s.taskRepo.Update(task)
}

func (s *taskService) RejectTask(taskID, approvalTaskID int64, req *model.ApprovalRequest, domainID, userID int64) error {
	// Get task
	task, err := s.taskRepo.GetByID(taskID, domainID)
	if err != nil {
		return err
	}

	// Get approval tasks
	approvalTasks, err := s.taskRepo.GetApprovalTasksByTaskID(taskID)
	if err != nil {
		return err
	}

	if len(approvalTasks) == 0 {
		return errors.New("no approval tasks found")
	}

	// Find the approval task to reject
	var currentApproval *model.ApprovalTask
	for i := range approvalTasks {
		if approvalTasks[i].ID == approvalTaskID {
			currentApproval = &approvalTasks[i]
			break
		}
	}

	if currentApproval == nil {
		return errors.New("approval task not found")
	}

	now := time.Now()

	// If sequence 1 is rejected, reject all sequences
	if currentApproval.Sequence == 1 {
		for i := range approvalTasks {
			approvalTasks[i].ApprovedBy = &userID
			approvalTasks[i].ApprovalStatusID = ptrInt64(helper.ApprovalStatusReject)
			approvalTasks[i].ApprovalDate = &now
			approvalTasks[i].Note = req.Note

			if err := s.taskRepo.UpdateApprovalTask(&approvalTasks[i]); err != nil {
				return err
			}
		}
	} else {
		// Reject only current sequence
		currentApproval.ApprovedBy = &userID
		currentApproval.ApprovalStatusID = ptrInt64(helper.ApprovalStatusReject)
		currentApproval.ApprovalDate = &now
		currentApproval.Note = req.Note

		if err := s.taskRepo.UpdateApprovalTask(currentApproval); err != nil {
			return err
		}
	}

	// Update task approval status to Reject
	approvalStatusID := int64(helper.ApprovalStatusReject)
	task.ApprovalStatusID = &approvalStatusID
	task.UpdatedBy = &userID

	return s.taskRepo.Update(task)
}

// Helper functions

func (s *taskService) uploadTaskFiles(taskID int64, files []*multipart.FileHeader, fileType int) error {
	var taskFiles []model.TaskFile

	for _, file := range files {
		// Upload file using SaveFile
		filePath, err := helper.SaveFile(file, helper.TaskFileUploadPath)
		if err != nil {
			return err
		}

		// Get file size
		fileSize := fmt.Sprintf("%d", file.Size)

		// Get content type
		contentType := file.Header.Get("Content-Type")

		taskFile := model.TaskFile{
			TaskID:       taskID,
			FileName:     file.Filename,
			FilePath:     filePath,
			FileSize:     &fileSize,
			FileType:     &contentType,
			TaskFileType: &fileType,
			Status:       true,
		}

		taskFiles = append(taskFiles, taskFile)
	}

	return s.taskRepo.CreateTaskFiles(taskFiles)
}

func (s *taskService) toTaskResponse(task *model.Task) *model.TaskResponse {
	resp := &model.TaskResponse{
		ID:                task.ID,
		DomainID:          task.DomainID,
		ProjectID:         task.ProjectID,
		Code:              task.Code,
		Title:             task.Title,
		Description:       task.Description,
		DescriptionBefore: task.DescriptionBefore,
		DescriptionAfter:  task.DescriptionAfter,
		Reason:            task.Reason,
		Revision:          task.Revision,
		Status:            task.Status,
		StatusID:          task.StatusID,
		PriorityID:        task.PriorityID,
		TypeID:            task.TypeID,
		StackID:           task.StackID,
		AssignedID:        task.AssignedID,
		CreatedBy:         task.CreatedBy,
		UpdatedBy:         task.UpdatedBy,
		ApprovedBy:        task.ApprovedBy,
		CompletedBy:       task.CompletedBy,
		DoneBy:            task.DoneBy,
		ApprovalStatusID:  task.ApprovalStatusID,
		StartDate:         task.StartDate,
		DueDate:           task.DueDate,
		CompletedDate:     task.CompletedDate,
		ApprovalDate:      task.ApprovalDate,
		DoneAt:            task.DoneAt,
		CreatedAt:         task.CreatedAt,
		UpdatedAt:         task.UpdatedAt,
	}

	// Convert relations
	if task.Project != nil {
		resp.Project = &model.ProjectResponse{
			ID:   task.Project.ID,
			Code: task.Project.Code,
			Name: task.Project.Name,
		}
	}

	if task.StatusTask != nil {
		resp.StatusTask = &model.ReferenceResponse{
			ID:   task.StatusTask.ID,
			Name: task.StatusTask.Name,
		}
	}

	if task.Priority != nil {
		resp.Priority = &model.ReferenceResponse{
			ID:   task.Priority.ID,
			Name: task.Priority.Name,
		}
	}

	if task.Type != nil {
		resp.Type = &model.ReferenceResponse{
			ID:   task.Type.ID,
			Name: task.Type.Name,
		}
	}

	if task.Stack != nil {
		resp.Stack = &model.ReferenceResponse{
			ID:   task.Stack.ID,
			Name: task.Stack.Name,
		}
	}

	if task.Assignee != nil {
		resp.Assignee = &model.UserResponse{
			ID:       task.Assignee.ID,
			Username: task.Assignee.Username,
			FullName: task.Assignee.FullName,
			Email:    task.Assignee.Email,
		}
	}

	if task.Creator != nil {
		resp.Creator = &model.UserResponse{
			ID:       task.Creator.ID,
			Username: task.Creator.Username,
			FullName: task.Creator.FullName,
			Email:    task.Creator.Email,
		}
	}

	if task.ApprovalStatus != nil {
		resp.ApprovalStatus = &model.ReferenceResponse{
			ID:   task.ApprovalStatus.ID,
			Name: task.ApprovalStatus.Name,
		}
	}

	// Convert task files
	if task.TaskFiles != nil {
		resp.TaskFiles = make([]model.TaskFileResponse, len(task.TaskFiles))
		for i, tf := range task.TaskFiles {
			resp.TaskFiles[i] = model.TaskFileResponse{
				ID:           tf.ID,
				TaskID:       tf.TaskID,
				FileName:     tf.FileName,
				FilePath:     tf.FilePath,
				FileSize:     tf.FileSize,
				FileType:     tf.FileType,
				TaskFileType: tf.TaskFileType,
				Status:       tf.Status,
				CreatedAt:    tf.CreatedAt,
				UpdatedAt:    tf.UpdatedAt,
			}
		}
	}

	// Convert approval tasks
	if task.ApprovalTasks != nil {
		resp.ApprovalTasks = make([]model.ApprovalTaskResponse, len(task.ApprovalTasks))
		for i, at := range task.ApprovalTasks {
			atResp := model.ApprovalTaskResponse{
				ID:               at.ID,
				TaskID:           at.TaskID,
				Sequence:         at.Sequence,
				ApprovedBy:       at.ApprovedBy,
				ApprovalStatusID: at.ApprovalStatusID,
				ApprovalDate:     at.ApprovalDate,
				Note:             at.Note,
				Status:           at.Status,
				CreatedAt:        at.CreatedAt,
				UpdatedAt:        at.UpdatedAt,
			}

			if at.Approver != nil {
				atResp.Approver = &model.UserResponse{
					ID:       at.Approver.ID,
					Username: at.Approver.Username,
					FullName: at.Approver.FullName,
					Email:    at.Approver.Email,
				}
			}

			if at.ApprovalStatus != nil {
				atResp.ApprovalStatus = &model.ReferenceResponse{
					ID:   at.ApprovalStatus.ID,
					Name: at.ApprovalStatus.Name,
				}
			}

			resp.ApprovalTasks[i] = atResp
		}
	}

	return resp
}

func ptrInt64(v int64) *int64 {
	return &v
}
