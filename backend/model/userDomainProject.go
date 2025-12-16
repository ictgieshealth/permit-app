package model

import (
	"time"
)

// UserDomainProject represents the junction table for user-domain-project relationship
type UserDomainProject struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"id"`
	UserID    int64     `gorm:"column:user_id;not null" json:"user_id"`
	DomainID  int64     `gorm:"column:domain_id;not null" json:"domain_id"`
	ProjectID int64     `gorm:"column:project_id;not null" json:"project_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relations
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Domain  *Domain  `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (UserDomainProject) TableName() string {
	return "user_domain_projects"
}
