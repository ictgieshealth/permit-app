package model

import "time"

type PermitType struct {
	ID                     int64     `json:"id" gorm:"primaryKey;column:id;autoIncrement"`
	DivisionID             *int64    `json:"division_id" gorm:"column:division_id"`
	Name                   string    `json:"name" gorm:"column:name;not null"`
	RiskPoint              *string   `json:"risk_point" gorm:"column:risk_point"`
	DefaultApplicationType *string   `json:"default_application_type" gorm:"column:default_application_type"`
	DefaultValidityPeriod  *string   `json:"default_validity_period" gorm:"column:default_validity_period"`
	Notes                  *string   `json:"notes" gorm:"column:notes;type:text"`
	CreatedAt              time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`

	Division *Division `json:"division,omitempty" gorm:"foreignKey:DivisionID;references:ID"`
}

func (PermitType) TableName() string {
	return "permit_types"
}

type PermitTypeRequest struct {
	DivisionID             *int64  `json:"division_id" form:"division_id"`
	Name                   string  `json:"name" form:"name" validate:"required"`
	RiskPoint              *string `json:"risk_point" form:"risk_point"`
	DefaultApplicationType *string `json:"default_application_type" form:"default_application_type"`
	DefaultValidityPeriod  *string `json:"default_validity_period" form:"default_validity_period"`
	Notes                  *string `json:"notes" form:"notes"`
}

type PermitTypeUpdateRequest struct {
	DivisionID             *int64  `json:"division_id" form:"division_id"`
	Name                   string  `json:"name" form:"name" validate:"required"`
	RiskPoint              *string `json:"risk_point" form:"risk_point"`
	DefaultApplicationType *string `json:"default_application_type" form:"default_application_type"`
	DefaultValidityPeriod  *string `json:"default_validity_period" form:"default_validity_period"`
	Notes                  *string `json:"notes" form:"notes"`
}

type PermitTypeResponse struct {
	ID                     int64             `json:"id"`
	DivisionID             *int64            `json:"division_id"`
	Name                   string            `json:"name"`
	RiskPoint              *string           `json:"risk_point"`
	DefaultApplicationType *string           `json:"default_application_type"`
	DefaultValidityPeriod  *string           `json:"default_validity_period"`
	Notes                  *string           `json:"notes"`
	CreatedAt              time.Time         `json:"created_at"`
	UpdatedAt              time.Time         `json:"updated_at"`
	Division               *DivisionResponse `json:"division,omitempty"`
}

type PermitTypeListRequest struct {
	DivisionID *int64 `json:"division_id" form:"division_id"`
	Name       string `json:"name" form:"name"`
	Page       int    `json:"page" form:"page" validate:"omitempty,min=1"`
	Limit      int    `json:"limit" form:"limit" validate:"omitempty,min=1,max=10000"`
}