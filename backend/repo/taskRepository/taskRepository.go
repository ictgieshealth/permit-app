package taskRepository

import (
	"permit-app/model"
	"time"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *model.Task) error
	GetByID(id int64, domainID int64) (*model.Task, error)
	GetByCode(code string, domainID int64) (*model.Task, error)
	GetAll(domainID int64, filters map[string]interface{}, page, limit int) ([]model.Task, int64, error)
	Update(task *model.Task) error
	Delete(id int64, domainID int64) error
	ChangeStatus(id int64, domainID int64, statusID int64, updatedBy int64) error
	ChangeType(id int64, domainID int64, typeID int64, updatedBy int64) error
	InReview(id int64, domainID int64, descBefore, descAfter *string, updatedBy int64) error
	SetReason(id int64, domainID int64, reason string, updatedBy int64) error
	SetRevision(id int64, domainID int64, revision *string, updatedBy int64) error
	GenerateCode(projectID int64) (string, error)

	// Approval related
	CreateApprovalTasks(approvalTasks []model.ApprovalTask) error
	GetApprovalTasksByTaskID(taskID int64) ([]model.ApprovalTask, error)
	GetApprovalTaskBySequence(taskID int64, sequence int16) (*model.ApprovalTask, error)
	UpdateApprovalTask(approvalTask *model.ApprovalTask) error
	UpdateTaskApprovalStatus(taskID, approvalStatusID int64, approvedBy *int64, approvalDate *time.Time, updatedBy int64) error

	// Task File related
	CreateTaskFiles(files []model.TaskFile) error
	GetTaskFilesByTaskID(taskID int64) ([]model.TaskFile, error)
	GetTaskFileByID(id int64) (*model.TaskFile, error)
	DeleteTaskFile(id int64) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) GetByID(id int64, domainID int64) (*model.Task, error) {
	var task model.Task
	err := r.db.Where("id = ? AND domain_id = ? AND deleted_at IS NULL", id, domainID).
		Preload("Domain").
		Preload("Project").
		Preload("StatusTask").
		Preload("Priority").
		Preload("Type").
		Preload("Stack").
		Preload("Assignee").
		Preload("Creator").
		Preload("Updater").
		Preload("Approver").
		Preload("Completer").
		Preload("DoneByUser").
		Preload("ApprovalStatus").
		Preload("TaskFiles").
		Preload("ApprovalTasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC").
				Preload("Approver").
				Preload("ApprovalStatus")
		}).
		First(&task).Error

	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetByCode(code string, domainID int64) (*model.Task, error) {
	var task model.Task
	err := r.db.Where("code = ? AND domain_id = ? AND deleted_at IS NULL", code, domainID).
		Preload("Domain").
		Preload("Project").
		Preload("StatusTask").
		Preload("Priority").
		Preload("Type").
		Preload("Stack").
		Preload("Assignee").
		Preload("Creator").
		Preload("ApprovalStatus").
		Preload("TaskFiles").
		Preload("ApprovalTasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC").
				Preload("Approver").
				Preload("ApprovalStatus")
		}).
		First(&task).Error

	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetAll(domainID int64, filters map[string]interface{}, page, limit int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	query := r.db.Model(&model.Task{}).Where("domain_id = ? AND deleted_at IS NULL", domainID)

	// Apply filters
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("LOWER(code) LIKE LOWER(?) OR LOWER(title) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%")
	}

	if projectID, ok := filters["project_id"].(int64); ok && projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}

	if statusID, ok := filters["status_id"].(int64); ok && statusID > 0 {
		query = query.Where("status_id = ?", statusID)
	}

	if approvalStatusID, ok := filters["approval_status_id"].(int64); ok && approvalStatusID > 0 {
		query = query.Where("approval_status_id = ?", approvalStatusID)
	}

	if assignedID, ok := filters["assigned_id"].(int64); ok && assignedID > 0 {
		query = query.Where("assigned_id = ?", assignedID)
	}

	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}

	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	// Get total count
	query.Count(&total)

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Preload("Project").
		Preload("StatusTask").
		Preload("Priority").
		Preload("Type").
		Preload("Stack").
		Preload("Assignee").
		Preload("Creator").
		Preload("ApprovalStatus").
		Preload("TaskFiles").
		Preload("ApprovalTasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC").
				Preload("Approver").
				Preload("ApprovalStatus")
		}).
		Find(&tasks).Error

	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) Update(task *model.Task) error {
	updates := make(map[string]interface{})
	
	// Always update these fields
	updates["title"] = task.Title
	updates["description"] = task.Description
	updates["description_before"] = task.DescriptionBefore
	updates["description_after"] = task.DescriptionAfter
	updates["reason"] = task.Reason
	updates["revision"] = task.Revision
	updates["status"] = task.Status
	updates["project_id"] = task.ProjectID
	updates["created_by"] = task.CreatedBy
	
	// Optional pointer fields - dereference to get actual values
	if task.StatusID != nil {
		updates["status_id"] = *task.StatusID
	}
	if task.PriorityID != nil {
		updates["priority_id"] = *task.PriorityID
	}
	if task.TypeID != nil {
		updates["type_id"] = *task.TypeID
	}
	if task.StackID != nil {
		updates["stack_id"] = *task.StackID
	}
	if task.AssignedID != nil {
		updates["assigned_id"] = *task.AssignedID
	}
	if task.UpdatedBy != nil {
		updates["updated_by"] = *task.UpdatedBy
	}
	if task.ApprovedBy != nil {
		updates["approved_by"] = *task.ApprovedBy
	}
	if task.CompletedBy != nil {
		updates["completed_by"] = *task.CompletedBy
	}
	if task.DoneBy != nil {
		updates["done_by"] = *task.DoneBy
	}
	if task.ApprovalStatusID != nil {
		updates["approval_status_id"] = *task.ApprovalStatusID
	}
	if task.StartDate != nil {
		updates["start_date"] = *task.StartDate
	}
	if task.DueDate != nil {
		updates["due_date"] = *task.DueDate
	}
	if task.CompletedDate != nil {
		updates["completed_date"] = *task.CompletedDate
	}
	if task.ApprovalDate != nil {
		updates["approval_date"] = *task.ApprovalDate
	}
	if task.DoneAt != nil {
		updates["done_at"] = *task.DoneAt
	}
	
	return r.db.Model(&model.Task{}).Where("id = ?", task.ID).Updates(updates).Error
}

func (r *taskRepository) Delete(id int64, domainID int64) error {
	now := time.Now()
	return r.db.Model(&model.Task{}).
		Where("id = ? AND domain_id = ?", id, domainID).
		Update("deleted_at", now).Error
}

func (r *taskRepository) ChangeStatus(id int64, domainID int64, statusID int64, updatedBy int64) error {
	updates := map[string]interface{}{
		"status_id":  statusID,
		"updated_by": updatedBy,
		"updated_at": time.Now(),
	}

	return r.db.Model(&model.Task{}).
		Where("id = ? AND domain_id = ? AND deleted_at IS NULL", id, domainID).
		Updates(updates).Error
}

func (r *taskRepository) ChangeType(id int64, domainID int64, typeID int64, updatedBy int64) error {
	updates := map[string]interface{}{
		"type_id":    typeID,
		"updated_by": updatedBy,
		"updated_at": time.Now(),
	}

	return r.db.Model(&model.Task{}).
		Where("id = ? AND domain_id = ? AND deleted_at IS NULL", id, domainID).
		Updates(updates).Error
}

func (r *taskRepository) InReview(id int64, domainID int64, descBefore, descAfter *string, updatedBy int64) error {
	updates := map[string]interface{}{
		"description_before": descBefore,
		"description_after":  descAfter,
		"updated_by":         updatedBy,
		"updated_at":         time.Now(),
	}

	return r.db.Model(&model.Task{}).
		Where("id = ? AND domain_id = ? AND deleted_at IS NULL", id, domainID).
		Updates(updates).Error
}

func (r *taskRepository) SetReason(id int64, domainID int64, reason string, updatedBy int64) error {
	updates := map[string]interface{}{
		"reason":     reason,
		"updated_by": updatedBy,
		"updated_at": time.Now(),
	}

	return r.db.Model(&model.Task{}).
		Where("id = ? AND domain_id = ? AND deleted_at IS NULL", id, domainID).
		Updates(updates).Error
}

func (r *taskRepository) SetRevision(id int64, domainID int64, revision *string, updatedBy int64) error {
	updates := map[string]interface{}{
		"revision":   revision,
		"updated_by": updatedBy,
		"updated_at": time.Now(),
	}

	return r.db.Model(&model.Task{}).
		Where("id = ? AND domain_id = ? AND deleted_at IS NULL", id, domainID).
		Updates(updates).Error
}

func (r *taskRepository) GenerateCode(projectID int64) (string, error) {
	var project model.Project
	err := r.db.First(&project, projectID).Error
	if err != nil {
		return "", err
	}

	var count int64
	r.db.Model(&model.Task{}).Where("project_id = ?", projectID).Count(&count)

	// Generate code format: PROJECT_CODE-TASK-XXXX
	code := project.Code + "-TASK-" + padLeft(int(count+1), 4)
	return code, nil
}

// Helper function to pad number with zeros
func padLeft(num, length int) string {
	str := ""
	for i := 0; i < length; i++ {
		str = "0" + str
	}
	numStr := ""
	for num > 0 {
		numStr = string(rune(num%10+'0')) + numStr
		num /= 10
	}
	if len(numStr) >= length {
		return numStr
	}
	return str[:length-len(numStr)] + numStr
}

// Approval Tasks
func (r *taskRepository) CreateApprovalTasks(approvalTasks []model.ApprovalTask) error {
	return r.db.Create(&approvalTasks).Error
}

func (r *taskRepository) GetApprovalTasksByTaskID(taskID int64) ([]model.ApprovalTask, error) {
	var approvalTasks []model.ApprovalTask
	err := r.db.Where("task_id = ?", taskID).
		Order("sequence ASC").
		Preload("Approver").
		Preload("ApprovalStatus").
		Find(&approvalTasks).Error

	return approvalTasks, err
}

func (r *taskRepository) GetApprovalTaskBySequence(taskID int64, sequence int16) (*model.ApprovalTask, error) {
	var approvalTask model.ApprovalTask
	err := r.db.Where("task_id = ? AND sequence = ?", taskID, sequence).
		Preload("Approver").
		Preload("ApprovalStatus").
		First(&approvalTask).Error

	if err != nil {
		return nil, err
	}
	return &approvalTask, nil
}

func (r *taskRepository) UpdateApprovalTask(approvalTask *model.ApprovalTask) error {
	updates := map[string]interface{}{
		"note": approvalTask.Note,
	}
	
	if approvalTask.ApprovedBy != nil {
		updates["approved_by"] = *approvalTask.ApprovedBy
	}
	if approvalTask.ApprovalStatusID != nil {
		updates["approval_status_id"] = *approvalTask.ApprovalStatusID
	}
	if approvalTask.ApprovalDate != nil {
		updates["approval_date"] = *approvalTask.ApprovalDate
	}
	
	return r.db.Model(&model.ApprovalTask{}).Where("id = ?", approvalTask.ID).Updates(updates).Error
}

// UpdateTaskApprovalStatus updates only the approval-related fields of a task
func (r *taskRepository) UpdateTaskApprovalStatus(taskID, approvalStatusID int64, approvedBy *int64, approvalDate *time.Time, updatedBy int64) error {
	updates := map[string]interface{}{
		"approval_status_id": approvalStatusID,
		"updated_by":         updatedBy,
	}
	
	if approvedBy != nil {
		updates["approved_by"] = *approvedBy
	}
	if approvalDate != nil {
		updates["approval_date"] = *approvalDate
	}
	
	return r.db.Model(&model.Task{}).Where("id = ?", taskID).Updates(updates).Error
}

// Task Files
func (r *taskRepository) CreateTaskFiles(files []model.TaskFile) error {
	if len(files) == 0 {
		return nil
	}
	return r.db.Create(&files).Error
}

func (r *taskRepository) GetTaskFilesByTaskID(taskID int64) ([]model.TaskFile, error) {
	var files []model.TaskFile
	err := r.db.Where("task_id = ?", taskID).Find(&files).Error
	return files, err
}

func (r *taskRepository) GetTaskFileByID(id int64) (*model.TaskFile, error) {
	var file model.TaskFile
	err := r.db.Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *taskRepository) DeleteTaskFile(id int64) error {
	return r.db.Delete(&model.TaskFile{}, id).Error
}
