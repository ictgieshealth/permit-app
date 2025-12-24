package model

import (
	"time"
)

// Project represents a project entity in the system
type Project struct {
	ID              int64      `gorm:"primaryKey;column:id" json:"id"`
	DomainID        int64      `gorm:"column:domain_id;not null" json:"domain_id"`
	Name            string     `gorm:"column:name;size:255;not null" json:"name"`
	Code            string     `gorm:"column:code;size:100;not null" json:"code"`
	Description     string     `gorm:"column:description;type:text" json:"description"`
	Status          bool       `gorm:"column:status;default:true" json:"status"`
	ProjectStatusID int64      `gorm:"column:project_status_id;not null" json:"project_status_id"`
	StartedAt       *time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt      *time.Time `gorm:"column:finished_at" json:"finished_at"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relations
	Domain        *Domain    `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
	ProjectStatus *Reference `gorm:"foreignKey:ProjectStatusID" json:"project_status,omitempty"`
	Users         []User     `gorm:"many2many:user_domain_projects;foreignKey:ID;joinForeignKey:ProjectID;References:ID;joinReferences:UserID" json:"users,omitempty"`
}

func (Project) TableName() string {
	return "projects"
}

// ProjectRequest represents the request body for creating/updating a project
type ProjectRequest struct {
	DomainID        int64   `json:"domain_id" validate:"required"`
	Name            string  `json:"name" validate:"required,max=255"`
	Code            string  `json:"code" validate:"omitempty,max=100"`
	Description     string  `json:"description" validate:"required"`
	Status          *bool   `json:"status"`
	ProjectStatusID *int64  `json:"project_status_id"`
	UserIDs         []int64 `json:"user_ids"`
}

// ProjectUpdateRequest represents the request body for updating a project
type ProjectUpdateRequest struct {
	Name            string  `json:"name" validate:"required,max=255"`
	Code            string  `json:"code" validate:"omitempty,max=100"`
	Description     string  `json:"description" validate:"required"`
	Status          *bool   `json:"status"`
	ProjectStatusID *int64  `json:"project_status_id"`
	UserIDs         []int64 `json:"user_ids"`
}

// ProjectStatusChangeRequest represents the request to change project status
type ProjectStatusChangeRequest struct {
	StatusID int64 `json:"status_id" validate:"required"`
}

// ProjectResponse represents the API response for a project
type ProjectResponse struct {
	ID              int64               `json:"id"`
	DomainID        int64               `json:"domain_id"`
	Name            string              `json:"name"`
	Code            string              `json:"code"`
	Description     string              `json:"description"`
	Status          bool                `json:"status"`
	ProjectStatusID int64               `json:"project_status_id"`
	StartedAt       *time.Time          `json:"started_at"`
	FinishedAt      *time.Time          `json:"finished_at"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	Domain          *DomainResponse     `json:"domain,omitempty"`
	ProjectStatus   *ReferenceResponse  `json:"project_status,omitempty"`
	Users           []UserBasicResponse `json:"users,omitempty"`
}

// UserBasicResponse for nested user data in project
type UserBasicResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

// ProjectListRequest represents the query parameters for listing projects
type ProjectListRequest struct {
	DomainID        int64  `form:"domain_id"`
	ProjectStatusID int64  `form:"project_status_id"`
	Name            string `form:"name"`
	Code            string `form:"code"`
	Status          *bool  `form:"status"`
	Page            int    `form:"page" validate:"omitempty,min=1"`
	Limit           int    `form:"limit" validate:"omitempty,min=1,max=10000"`
}

// ToProjectResponse converts Project entity to ProjectResponse
func ToProjectResponse(project *Project) ProjectResponse {
	response := ProjectResponse{
		ID:              project.ID,
		DomainID:        project.DomainID,
		Name:            project.Name,
		Code:            project.Code,
		Description:     project.Description,
		Status:          project.Status,
		ProjectStatusID: project.ProjectStatusID,
		StartedAt:       project.StartedAt,
		FinishedAt:      project.FinishedAt,
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
	}

	if project.Domain != nil {
		response.Domain = &DomainResponse{
			ID:          project.Domain.ID,
			Code:        project.Domain.Code,
			Name:        project.Domain.Name,
			Description: project.Domain.Description,
			IsActive:    project.Domain.IsActive,
			CreatedAt:   project.Domain.CreatedAt,
			UpdatedAt:   project.Domain.UpdatedAt,
		}
	}

	if project.ProjectStatus != nil {
		statusResp := ToReferenceResponse(project.ProjectStatus)
		response.ProjectStatus = &statusResp
	}

	if len(project.Users) > 0 {
		response.Users = make([]UserBasicResponse, len(project.Users))
		for i, user := range project.Users {
			response.Users[i] = UserBasicResponse{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
				FullName: user.FullName,
			}
		}
	}

	return response
}
