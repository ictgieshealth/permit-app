package model

import (
	"time"
)

type Task struct {
	ID                int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	DomainID          int64      `gorm:"not null;index" json:"domain_id"`
	ProjectID         int64      `gorm:"not null;index" json:"project_id"`
	Code              string     `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Title             string     `gorm:"size:255;not null" json:"title"`
	Description       *string    `gorm:"type:text" json:"description"`
	DescriptionBefore *string    `gorm:"type:text;column:description_before" json:"description_before"`
	DescriptionAfter  *string    `gorm:"type:text;column:description_after" json:"description_after"`
	Reason            *string    `gorm:"type:text" json:"reason"`
	Revision          *string    `gorm:"type:text" json:"revision"`
	Status            bool       `gorm:"default:true" json:"status"`
	StatusID          *int64     `gorm:"index" json:"status_id"`
	PriorityID        *int64     `gorm:"index" json:"priority_id"`
	TypeID            *int64     `gorm:"index" json:"type_id"`
	StackID           *int64     `gorm:"index" json:"stack_id"`
	AssignedID        *int64     `gorm:"index" json:"assigned_id"`
	CreatedBy         int64      `gorm:"not null;index" json:"created_by"`
	UpdatedBy         *int64     `json:"updated_by"`
	ApprovedBy        *int64     `json:"approved_by"`
	CompletedBy       *int64     `json:"completed_by"`
	DoneBy            *int64     `json:"done_by"`
	ApprovalStatusID  *int64     `gorm:"index" json:"approval_status_id"`
	StartDate         *time.Time `json:"start_date"`
	DueDate           *time.Time `json:"due_date"`
	CompletedDate     *time.Time `json:"completed_date"`
	ApprovalDate      *time.Time `json:"approval_date"`
	DoneAt            *time.Time `json:"done_at"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Domain         *Domain        `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
	Project        *Project       `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	StatusTask     *Reference     `gorm:"foreignKey:StatusID" json:"status_task,omitempty"`
	Priority       *Reference     `gorm:"foreignKey:PriorityID" json:"priority,omitempty"`
	Type           *Reference     `gorm:"foreignKey:TypeID" json:"type,omitempty"`
	Stack          *Reference     `gorm:"foreignKey:StackID" json:"stack,omitempty"`
	Assignee       *User          `gorm:"foreignKey:AssignedID" json:"assignee,omitempty"`
	Creator        *User          `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	Updater        *User          `gorm:"foreignKey:UpdatedBy" json:"updated_by_user,omitempty"`
	Approver       *User          `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
	Completer      *User          `gorm:"foreignKey:CompletedBy" json:"completed_by_user,omitempty"`
	DoneByUser     *User          `gorm:"foreignKey:DoneBy" json:"done_by_user,omitempty"`
	ApprovalStatus *Reference     `gorm:"foreignKey:ApprovalStatusID" json:"approval_status,omitempty"`
	TaskFiles      []TaskFile     `gorm:"foreignKey:TaskID" json:"task_files,omitempty"`
	ApprovalTasks  []ApprovalTask `gorm:"foreignKey:TaskID" json:"approval_tasks,omitempty"`
}

type TaskFile struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID       int64     `gorm:"not null;index" json:"task_id"`
	FileName     string    `gorm:"size:255;not null" json:"file_name"`
	FilePath     string    `gorm:"size:500;not null" json:"file_path"`
	FileSize     *string   `gorm:"size:50" json:"file_size"`
	FileType     *string   `gorm:"size:100" json:"file_type"`
	TaskFileType *int      `json:"task_file_type"`
	Status       bool      `gorm:"default:true" json:"status"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Task *Task `gorm:"foreignKey:TaskID" json:"task,omitempty"`
}

type ApprovalTask struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID           int64      `gorm:"not null;index" json:"task_id"`
	Sequence         int16      `gorm:"not null" json:"sequence"`
	ApprovedBy       *int64     `gorm:"index" json:"approved_by"`
	ApprovalStatusID *int64     `gorm:"index" json:"approval_status_id"`
	ApprovalDate     *time.Time `json:"approval_date"`
	Note             *string    `gorm:"size:500" json:"note"`
	Status           bool       `gorm:"default:true" json:"status"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Task           *Task      `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Approver       *User      `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
	ApprovalStatus *Reference `gorm:"foreignKey:ApprovalStatusID" json:"approval_status,omitempty"`
}

func (Task) TableName() string {
	return "tasks"
}

func (TaskFile) TableName() string {
	return "task_files"
}

func (ApprovalTask) TableName() string {
	return "approval_tasks"
}

// Request & Response DTOs

type TaskRequest struct {
	ProjectID   int64   `json:"project_id" validate:"required"`
	Title       string  `json:"title" validate:"required,max=255"`
	Description *string `json:"description"`
	PriorityID  int64   `json:"priority_id" validate:"required"`
	AssignedID  *int64  `json:"assigned_id"`
	StackID     *int64  `json:"stack_id"`
	DueDate     *string `json:"due_date"`
}

type TaskUpdateRequest struct {
	Title             string  `json:"title" validate:"required,max=255"`
	Description       *string `json:"description"`
	DescriptionBefore *string `json:"description_before"`
	DescriptionAfter  *string `json:"description_after"`
	ProjectID         int64   `json:"project_id" validate:"required"`
	PriorityID        int64   `json:"priority_id" validate:"required"`
	AssignedID        *int64  `json:"assigned_id"`
	StackID           *int64  `json:"stack_id"`
	DueDate           *string `json:"due_date"`
}

type TaskChangeStatusRequest struct {
	StatusID int64 `json:"status_id" validate:"required"`
}

type TaskChangeTypeRequest struct {
	TypeID int64 `json:"type_id" validate:"required"`
}

type TaskInReviewRequest struct {
	DescriptionBefore *string `json:"description_before"`
	DescriptionAfter  *string `json:"description_after"`
}

type TaskReasonRequest struct {
	Reason string `json:"reason" validate:"required"`
}

type TaskRevisionRequest struct {
	Revision *string `json:"revision"`
}

type ApprovalRequest struct {
	Note *string `json:"note"`
}

type TaskListRequest struct {
	Search           *string `form:"search"`
	ProjectID        *int64  `form:"project_id"`
	StatusID         *int64  `form:"status_id"`
	ApprovalStatusID *int64  `form:"approval_status_id"`
	AssignedID       *int64  `form:"assigned_id"`
	StartDate        *string `form:"start_date"`
	EndDate          *string `form:"end_date"`
	Page             int     `form:"page" validate:"min=1"`
	Limit            int     `form:"limit" validate:"min=1,max=100"`
}

type TaskResponse struct {
	ID                int64                  `json:"id"`
	DomainID          int64                  `json:"domain_id"`
	ProjectID         int64                  `json:"project_id"`
	Code              string                 `json:"code"`
	Title             string                 `json:"title"`
	Description       *string                `json:"description"`
	DescriptionBefore *string                `json:"description_before"`
	DescriptionAfter  *string                `json:"description_after"`
	Reason            *string                `json:"reason"`
	Revision          *string                `json:"revision"`
	Status            bool                   `json:"status"`
	StatusID          *int64                 `json:"status_id"`
	PriorityID        *int64                 `json:"priority_id"`
	TypeID            *int64                 `json:"type_id"`
	StackID           *int64                 `json:"stack_id"`
	AssignedID        *int64                 `json:"assigned_id"`
	CreatedBy         int64                  `json:"created_by"`
	UpdatedBy         *int64                 `json:"updated_by"`
	ApprovedBy        *int64                 `json:"approved_by"`
	CompletedBy       *int64                 `json:"completed_by"`
	DoneBy            *int64                 `json:"done_by"`
	ApprovalStatusID  *int64                 `json:"approval_status_id"`
	StartDate         *time.Time             `json:"start_date"`
	DueDate           *time.Time             `json:"due_date"`
	CompletedDate     *time.Time             `json:"completed_date"`
	ApprovalDate      *time.Time             `json:"approval_date"`
	DoneAt            *time.Time             `json:"done_at"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	Project           *ProjectResponse       `json:"project,omitempty"`
	StatusTask        *ReferenceResponse     `json:"status_task,omitempty"`
	Priority          *ReferenceResponse     `json:"priority,omitempty"`
	Type              *ReferenceResponse     `json:"type,omitempty"`
	Stack             *ReferenceResponse     `json:"stack,omitempty"`
	Assignee          *UserResponse          `json:"assignee,omitempty"`
	Creator           *UserResponse          `json:"created_by_user,omitempty"`
	ApprovalStatus    *ReferenceResponse     `json:"approval_status,omitempty"`
	TaskFiles         []TaskFileResponse     `json:"task_files,omitempty"`
	ApprovalTasks     []ApprovalTaskResponse `json:"approval_tasks,omitempty"`
}

type TaskFileResponse struct {
	ID           int64     `json:"id"`
	TaskID       int64     `json:"task_id"`
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	FileSize     *string   `json:"file_size"`
	FileType     *string   `json:"file_type"`
	TaskFileType *int      `json:"task_file_type"`
	Status       bool      `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ApprovalTaskResponse struct {
	ID               int64              `json:"id"`
	TaskID           int64              `json:"task_id"`
	Sequence         int16              `json:"sequence"`
	ApprovedBy       *int64             `json:"approved_by"`
	ApprovalStatusID *int64             `json:"approval_status_id"`
	ApprovalDate     *time.Time         `json:"approval_date"`
	Note             *string            `json:"note"`
	Status           bool               `json:"status"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	Approver         *UserResponse      `json:"approved_by_user,omitempty"`
	ApprovalStatus   *ReferenceResponse `json:"approval_status,omitempty"`
}
