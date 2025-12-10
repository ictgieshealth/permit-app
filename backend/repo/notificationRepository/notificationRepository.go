package notificationRepository

import (
	"permit-app/model"
	"time"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *model.Notification) error
	FindByUserID(userID int64, limit int, offset int) ([]model.Notification, error)
	FindUnreadByUserID(userID int64) ([]model.Notification, error)
	CountUnreadByUserID(userID int64) (int64, error)
	MarkAsRead(notificationIDs []int64, userID int64) error
	MarkAllAsRead(userID int64) error
	FindByID(id int64) (*model.Notification, error)
	Delete(id int64) error
	CheckExistingNotification(permitID int64, notificationType string, createdAfter time.Time) (bool, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindByUserID(userID int64, limit int, offset int) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.Where("user_id = ?", userID).
		Preload("Permit").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) FindUnreadByUserID(userID int64) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.Where("user_id = ? AND is_read = ?", userID, false).
		Preload("Permit").
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) CountUnreadByUserID(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) MarkAsRead(notificationIDs []int64, userID int64) error {
	now := time.Now()
	return r.db.Model(&model.Notification{}).
		Where("id IN ? AND user_id = ?", notificationIDs, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *notificationRepository) MarkAllAsRead(userID int64) error {
	now := time.Now()
	return r.db.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *notificationRepository) FindByID(id int64) (*model.Notification, error) {
	var notification model.Notification
	err := r.db.Preload("Permit").First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) Delete(id int64) error {
	return r.db.Delete(&model.Notification{}, id).Error
}

func (r *notificationRepository) CheckExistingNotification(permitID int64, notificationType string, createdAfter time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).
		Where("permit_id = ? AND type = ? AND created_at >= ?", permitID, notificationType, createdAfter).
		Count(&count).Error
	return count > 0, err
}
