package model

import "time"

type Module struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"id"`
	Code      string    `gorm:"unique;not null;column:code" json:"code"`
	Name      string    `gorm:"not null;column:name" json:"name"`
	IsActive  bool      `gorm:"not null;default:true;column:is_active" json:"is_active"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Module) TableName() string {
	return "modules"
}

// Request/Response DTOs
type ModuleRequest struct {
	Code     string `json:"code" validate:"required,max=50"`
	Name     string `json:"name" validate:"required,max=100"`
	IsActive *bool  `json:"is_active"`
}

type ModuleResponse struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ModuleListRequest struct {
	Code     string `form:"code"`
	Name     string `form:"name"`
	IsActive *bool  `form:"is_active"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}

func ToModuleResponse(module *Module) ModuleResponse {
	return ModuleResponse{
		ID:        module.ID,
		Code:      module.Code,
		Name:      module.Name,
		IsActive:  module.IsActive,
		CreatedAt: module.CreatedAt,
		UpdatedAt: module.UpdatedAt,
	}
}
