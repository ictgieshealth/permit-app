package model

import "time"

type Role struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type RoleRequest struct {
	Name string `json:"name" validate:"required,max=50"`
}

type RoleUpdateRequest struct {
	Name string `json:"name" validate:"required,max=50"`
}

type RoleResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleListRequest struct {
	Name  string `form:"name"`
	Page  int    `form:"page" validate:"omitempty,min=1"`
	Limit int    `form:"limit" validate:"omitempty,min=1,max=10000"`
}
