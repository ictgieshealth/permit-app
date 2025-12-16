package model

import "time"

type Reference struct {
	ID                  int64     `gorm:"primaryKey;column:id" json:"id"`
	ReferenceCategoryID int64     `gorm:"not null;column:reference_category_id" json:"reference_category_id"`
	Name                string    `gorm:"not null;column:name" json:"name"`
	IsActive            bool      `gorm:"not null;default:true;column:is_active" json:"is_active"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at" json:"updated_at"`

	// Relations
	ReferenceCategory *ReferenceCategory `gorm:"foreignKey:ReferenceCategoryID" json:"reference_category,omitempty"`
}

func (Reference) TableName() string {
	return "references"
}

// Request/Response DTOs
type ReferenceRequest struct {
	ReferenceCategoryID int64  `json:"reference_category_id" validate:"required"`
	Name                string `json:"name" validate:"required,max=100"`
	IsActive            *bool  `json:"is_active"`
}

type ReferenceResponse struct {
	ID                  int64                      `json:"id"`
	ReferenceCategoryID int64                      `json:"reference_category_id"`
	Name                string                     `json:"name"`
	IsActive            bool                       `json:"is_active"`
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
	ReferenceCategory   *ReferenceCategoryResponse `json:"reference_category,omitempty"`
}

type ReferenceListRequest struct {
	ReferenceCategoryID int64  `form:"reference_category_id"`
	ModuleID            int64  `form:"module_id"`
	Name                string `form:"name"`
	IsActive            *bool  `form:"is_active"`
	Page                int    `form:"page"`
	Limit               int    `form:"limit"`
}

func ToReferenceResponse(reference *Reference) ReferenceResponse {
	response := ReferenceResponse{
		ID:                  reference.ID,
		ReferenceCategoryID: reference.ReferenceCategoryID,
		Name:                reference.Name,
		IsActive:            reference.IsActive,
		CreatedAt:           reference.CreatedAt,
		UpdatedAt:           reference.UpdatedAt,
	}

	if reference.ReferenceCategory != nil {
		categoryResp := ToReferenceCategoryResponse(reference.ReferenceCategory)
		response.ReferenceCategory = &categoryResp
	}

	return response
}
