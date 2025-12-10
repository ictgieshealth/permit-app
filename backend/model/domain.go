package model

import "time"

type Domain struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type DomainRequest struct {
	Code        string  `json:"code" validate:"required,max=50"`
	Name        string  `json:"name" validate:"required,max=255"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type DomainUpdateRequest struct {
	Code        string  `json:"code" validate:"omitempty,max=50"`
	Name        string  `json:"name" validate:"omitempty,max=255"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type DomainResponse struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DomainListRequest struct {
	Code     string `form:"code"`
	Name     string `form:"name"`
	IsActive *bool  `form:"is_active"`
	Page     int    `form:"page" validate:"omitempty,min=1"`
	Limit    int    `form:"limit" validate:"omitempty,min=1,max=10000"`
}
