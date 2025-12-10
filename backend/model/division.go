package model

import "time"

type Division struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DomainID  int64     `gorm:"not null" json:"domain_id"`
	Code      string    `gorm:"type:varchar(50);not null" json:"code"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Domain *Domain `json:"domain,omitempty" gorm:"foreignKey:DomainID;references:ID"`
}

type DivisionRequest struct {
	DomainID int64  `json:"domain_id" validate:"required"`
	Code     string `json:"code" validate:"required,max=50"`
	Name     string `json:"name" validate:"required,max=255"`
}

type DivisionUpdateRequest struct {
	DomainID int64  `json:"domain_id" validate:"omitempty"`
	Code     string `json:"code" validate:"omitempty,max=50"`
	Name     string `json:"name" validate:"omitempty,max=255"`
}

type DivisionResponse struct {
	ID        int64           `json:"id"`
	DomainID  int64           `json:"domain_id"`
	Code      string          `json:"code"`
	Name      string          `json:"name"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Domain    *DomainResponse `json:"domain,omitempty"`
}

type DivisionListRequest struct {
	DomainID *int64 `form:"domain_id"`
	Code     string `form:"code"`
	Name     string `form:"name"`
	Page     int    `form:"page" validate:"omitempty,min=1"`
	Limit    int    `form:"limit" validate:"omitempty,min=1,max=10000"`
}
