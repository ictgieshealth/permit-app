package model

import "time"

type Role struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"type:varchar(50);not null" json:"name"`
	Category    string    `gorm:"type:varchar(50);not null" json:"category"` // 'Permit' or 'Ticketing'
	Description *string   `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type RoleRequest struct {
	Code        string  `json:"code" validate:"required,max=50"`
	Name        string  `json:"name" validate:"required,max=50"`
	Category    string  `json:"category" validate:"required,oneof=Permit Ticketing"`
	Description *string `json:"description"`
}

type RoleUpdateRequest struct {
	Code        string  `json:"code" validate:"omitempty,max=50"`
	Name        string  `json:"name" validate:"omitempty,max=50"`
	Category    string  `json:"category" validate:"omitempty,oneof=Permit Ticketing"`
	Description *string `json:"description"`
}

type RoleResponse struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoleListRequest struct {
	Code     string `form:"code"`
	Name     string `form:"name"`
	Category string `form:"category" validate:"omitempty,oneof=Permit Ticketing"`
	Page     int    `form:"page" validate:"omitempty,min=1"`
	Limit    int    `form:"limit" validate:"omitempty,min=1,max=10000"`
}
