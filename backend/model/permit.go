package model

import (
	"permit-app/helper"
	"time"
)

type Permit struct {
	ID                     int64     `json:"id" gorm:"primaryKey;column:id;autoIncrement"`
	DomainID               int64     `json:"domain_id" gorm:"column:domain_id;not null"`
	DivisionID             *int64    `json:"division_id" gorm:"column:division_id"`
	PermitTypeID           int64     `json:"permit_type_id" gorm:"column:permit_type_id;not null"`
	Name                   string    `json:"name" gorm:"column:name;not null"`
	ApplicationType        string    `json:"application_type" gorm:"column:application_type;not null"`
	PermitNo               string    `json:"permit_no" gorm:"column:permit_no;not null"`
	EffectiveDate          time.Time `json:"effective_date" gorm:"column:effective_date;not null"`
	ExpiryDate             time.Time `json:"expiry_date" gorm:"column:expiry_date;not null"`
	EffectiveTerm          *string   `json:"effective_term" gorm:"column:effective_term"`
	ResponsiblePersonID    *int64    `json:"responsible_person_id" gorm:"column:responsible_person_id"`
	ResponsibleDocPersonID *int64    `json:"responsible_doc_person_id" gorm:"column:responsible_doc_person_id"`
	DocName                *string   `json:"doc_name" gorm:"column:doc_name"`
	DocNumber              *string   `json:"doc_number" gorm:"column:doc_number"`
	DocFileName            *string   `json:"doc_file_name" gorm:"column:doc_file_name"`
	DocFilePath            *string   `json:"doc_file_path" gorm:"column:doc_file_path"`
	DocFileSize            *int64    `json:"doc_file_size" gorm:"column:doc_file_size"`
	DocFileType            *string   `json:"doc_file_type" gorm:"column:doc_file_type"`
	Status                 string    `json:"status" gorm:"column:status;not null;default:'active'"`
	CreatedAt              time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`

	Domain              *Domain     `json:"domain,omitempty" gorm:"foreignKey:DomainID;references:ID"`
	Division            *Division   `json:"division,omitempty" gorm:"foreignKey:DivisionID;references:ID"`
	PermitType          *PermitType `json:"permit_type,omitempty" gorm:"foreignKey:PermitTypeID;references:ID"`
	ResponsiblePerson   *User       `json:"responsible_person,omitempty" gorm:"foreignKey:ResponsiblePersonID;references:ID"`
	ResponsibleDocPerson *User      `json:"responsible_doc_person,omitempty" gorm:"foreignKey:ResponsibleDocPersonID;references:ID"`
}

func (Permit) TableName() string {
	return "permits"
}

// PermitFormRequest is for multipart/form-data binding
type PermitFormRequest struct {
	DomainID               int64   `form:"domain_id" validate:"required"`
	DivisionID             *int64  `form:"division_id"`
	PermitTypeID           int64   `form:"permit_type_id" validate:"required"`
	Name                   string  `form:"name" validate:"required"`
	ApplicationType        string  `form:"application_type" validate:"required"`
	PermitNo               string  `form:"permit_no" validate:"required"`
	EffectiveDate          string  `form:"effective_date" validate:"required"`
	ExpiryDate             string  `form:"expiry_date" validate:"required"`
	EffectiveTerm          *string `form:"effective_term"`
	ResponsiblePersonID    *int64  `form:"responsible_person_id"`
	ResponsibleDocPersonID *int64  `form:"responsible_doc_person_id"`
	DocName                *string `form:"doc_name"`
	DocNumber              *string `form:"doc_number"`
	Status                 string  `form:"status" validate:"required"`
}

type PermitRequest struct {
	DomainID               int64       `json:"domain_id" form:"domain_id" validate:"required"`
	DivisionID             *int64      `json:"division_id" form:"division_id"`
	PermitTypeID           int64       `json:"permit_type_id" form:"permit_type_id" validate:"required"`
	Name                   string      `json:"name" form:"name" validate:"required"`
	ApplicationType        string      `json:"application_type" form:"application_type" validate:"required"`
	PermitNo               string      `json:"permit_no" form:"permit_no" validate:"required"`
	EffectiveDate          helper.Date `json:"effective_date" form:"effective_date" validate:"required"`
	ExpiryDate             helper.Date `json:"expiry_date" form:"expiry_date" validate:"required"`
	EffectiveTerm          *string     `json:"effective_term" form:"effective_term"`
	ResponsiblePersonID    *int64      `json:"responsible_person_id" form:"responsible_person_id"`
	ResponsibleDocPersonID *int64      `json:"responsible_doc_person_id" form:"responsible_doc_person_id"`
	DocName                *string     `json:"doc_name" form:"doc_name"`
	DocNumber              *string     `json:"doc_number" form:"doc_number"`
	DocFileName            *string     `json:"doc_file_name" form:"doc_file_name"`
	DocFilePath            *string     `json:"doc_file_path" form:"doc_file_path"`
	DocFileSize            *int64      `json:"doc_file_size" form:"doc_file_size"`
	DocFileType            *string     `json:"doc_file_type" form:"doc_file_type"`
	Status                 string      `json:"status" form:"status" validate:"required"`
}

// PermitFormUpdateRequest is for multipart/form-data binding
type PermitFormUpdateRequest struct {
	DomainID               int64   `form:"domain_id" validate:"required"`
	DivisionID             *int64  `form:"division_id"`
	PermitTypeID           int64   `form:"permit_type_id" validate:"required"`
	Name                   string  `form:"name" validate:"required"`
	ApplicationType        string  `form:"application_type" validate:"required"`
	PermitNo               string  `form:"permit_no" validate:"required"`
	EffectiveDate          string  `form:"effective_date" validate:"required"`
	ExpiryDate             string  `form:"expiry_date" validate:"required"`
	EffectiveTerm          *string `form:"effective_term"`
	ResponsiblePersonID    *int64  `form:"responsible_person_id"`
	ResponsibleDocPersonID *int64  `form:"responsible_doc_person_id"`
	DocName                *string `form:"doc_name"`
	DocNumber              *string `form:"doc_number"`
	Status                 string  `form:"status" validate:"required"`
}

type PermitUpdateRequest struct {
	DomainID               int64       `json:"domain_id" form:"domain_id" validate:"required"`
	DivisionID             *int64      `json:"division_id" form:"division_id"`
	PermitTypeID           int64       `json:"permit_type_id" form:"permit_type_id" validate:"required"`
	Name                   string      `json:"name" form:"name" validate:"required"`
	ApplicationType        string      `json:"application_type" form:"application_type" validate:"required"`
	PermitNo               string      `json:"permit_no" form:"permit_no" validate:"required"`
	EffectiveDate          helper.Date `json:"effective_date" form:"effective_date" validate:"required"`
	ExpiryDate             helper.Date `json:"expiry_date" form:"expiry_date" validate:"required"`
	EffectiveTerm          *string     `json:"effective_term" form:"effective_term"`
	ResponsiblePersonID    *int64      `json:"responsible_person_id" form:"responsible_person_id"`
	ResponsibleDocPersonID *int64      `json:"responsible_doc_person_id" form:"responsible_doc_person_id"`
	DocName                *string   `json:"doc_name" form:"doc_name"`
	DocNumber              *string   `json:"doc_number" form:"doc_number"`
	DocFileName            *string   `json:"doc_file_name" form:"doc_file_name"`
	DocFilePath            *string   `json:"doc_file_path" form:"doc_file_path"`
	DocFileSize            *int64    `json:"doc_file_size" form:"doc_file_size"`
	DocFileType            *string   `json:"doc_file_type" form:"doc_file_type"`
	Status                 string    `json:"status" form:"status" validate:"required"`
}

type PermitResponse struct {
	ID                     int64               `json:"id"`
	DomainID               int64               `json:"domain_id"`
	DivisionID             *int64              `json:"division_id"`
	PermitTypeID           int64               `json:"permit_type_id"`
	Name                   string              `json:"name"`
	ApplicationType        string              `json:"application_type"`
	PermitNo               string              `json:"permit_no"`
	EffectiveDate          time.Time           `json:"effective_date"`
	ExpiryDate             time.Time           `json:"expiry_date"`
	EffectiveTerm          *string             `json:"effective_term"`
	ResponsiblePersonID    *int64              `json:"responsible_person_id"`
	ResponsibleDocPersonID *int64              `json:"responsible_doc_person_id"`
	DocName                *string             `json:"doc_name"`
	DocNumber              *string             `json:"doc_number"`
	DocFileName            *string             `json:"doc_file_name"`
	DocFilePath            *string             `json:"doc_file_path"`
	DocFileSize            *int64              `json:"doc_file_size"`
	DocFileType            *string             `json:"doc_file_type"`
	Status                 string              `json:"status"`
	CreatedAt              time.Time           `json:"created_at"`
	UpdatedAt              time.Time           `json:"updated_at"`
	Domain                 *DomainResponse     `json:"domain,omitempty"`
	Division               *DivisionResponse   `json:"division,omitempty"`
	PermitType             *PermitTypeResponse `json:"permit_type,omitempty"`
	ResponsiblePerson      *UserResponse       `json:"responsible_person,omitempty"`
	ResponsibleDocPerson   *UserResponse       `json:"responsible_doc_person,omitempty"`
}

type PermitListRequest struct {
	DomainID          *int64 `json:"domain_id" form:"domain_id"`
	DivisionID        *int64 `json:"division_id" form:"division_id"`
	PermitTypeID      *int64 `json:"permit_type_id" form:"permit_type_id"`
	Name              string `json:"name" form:"name"`
	ApplicationType   string `json:"application_type" form:"application_type"`
	PermitNo          string `json:"permit_no" form:"permit_no"`
	ResponsiblePerson string `json:"responsible_person" form:"responsible_person"`
	Status            string `json:"status" form:"status"`
	Page              int    `json:"page" form:"page" validate:"omitempty,min=1"`
	Limit             int    `json:"limit" form:"limit" validate:"omitempty,min=1,max=10000"`
}
