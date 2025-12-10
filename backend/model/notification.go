package model

import (
	"time"
)

type Notification struct {
	ID        int64     `json:"id" gorm:"primaryKey;column:id;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"column:user_id;not null"`
	PermitID  int64     `json:"permit_id" gorm:"column:permit_id;not null"`
	Type      string    `json:"type" gorm:"column:type;not null"` // expiry_reminder, expiry_warning, expired
	Title     string    `json:"title" gorm:"column:title;not null"`
	Message   string    `json:"message" gorm:"column:message;not null"`
	IsRead    bool      `json:"is_read" gorm:"column:is_read;default:false"`
	ReadAt    *time.Time `json:"read_at" gorm:"column:read_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`

	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Permit *Permit `json:"permit,omitempty" gorm:"foreignKey:PermitID;references:ID"`
}

func (Notification) TableName() string {
	return "notifications"
}

type NotificationResponse struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	PermitID  int64      `json:"permit_id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
	Permit    *struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		PermitNo  string `json:"permit_no"`
		ExpiryDate time.Time `json:"expiry_date"`
	} `json:"permit,omitempty"`
}

type MarkAsReadRequest struct {
	NotificationIDs []int64 `json:"notification_ids" validate:"required"`
}
