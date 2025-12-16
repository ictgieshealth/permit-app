package model

import "time"

type ReferenceCategory struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"id"`
	ModuleID  int64     `gorm:"not null;column:module_id" json:"module_id"`
	Name      string    `gorm:"not null;column:name" json:"name"`
	IsActive  bool      `gorm:"not null;default:true;column:is_active" json:"is_active"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// Relations
	Module     *Module      `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	References []*Reference `gorm:"foreignKey:ReferenceCategoryID" json:"references,omitempty"`
}

func (ReferenceCategory) TableName() string {
	return "reference_categories"
}

// Request/Response DTOs
type ReferenceCategoryRequest struct {
	ModuleID int64  `json:"module_id" validate:"required"`
	Name     string `json:"name" validate:"required,max=100"`
	IsActive *bool  `json:"is_active"`
}

type ReferenceCategoryResponse struct {
	ID        int64           `json:"id"`
	ModuleID  int64           `json:"module_id"`
	Name      string          `json:"name"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Module    *ModuleResponse `json:"module,omitempty"`
}

type ReferenceCategoryListRequest struct {
	ModuleID int64  `form:"module_id"`
	Name     string `form:"name"`
	IsActive *bool  `form:"is_active"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}

func ToReferenceCategoryResponse(category *ReferenceCategory) ReferenceCategoryResponse {
	response := ReferenceCategoryResponse{
		ID:        category.ID,
		ModuleID:  category.ModuleID,
		Name:      category.Name,
		IsActive:  category.IsActive,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	if category.Module != nil {
		moduleResp := ToModuleResponse(category.Module)
		response.Module = &moduleResp
	}

	return response
}
