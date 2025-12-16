package notificationService

import (
	"errors"
	"fmt"
	"permit-app/helper"
	"permit-app/model"
	"permit-app/repo/notificationRepository"
	"permit-app/repo/permitRepository"
	"permit-app/repo/userRepository"
	"time"
)

type NotificationService interface {
	GetUserNotifications(userID int64, limit int, offset int) ([]model.NotificationResponse, error)
	GetUnreadNotifications(userID int64) ([]model.NotificationResponse, error)
	GetUnreadCount(userID int64) (int64, error)
	MarkAsRead(notificationIDs []int64, userID int64) error
	MarkAllAsRead(userID int64) error
	DeleteNotification(id int64, userID int64) error
	CheckAndSendExpiryNotifications() error
}

type notificationService struct {
	notificationRepo notificationRepository.NotificationRepository
	permitRepo       permitRepository.PermitRepository
	userRepo         userRepository.UserRepository
}

func NewNotificationService(
	notificationRepo notificationRepository.NotificationRepository,
	permitRepo permitRepository.PermitRepository,
	userRepo userRepository.UserRepository,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		permitRepo:       permitRepo,
		userRepo:         userRepo,
	}
}

func (s *notificationService) GetUserNotifications(userID int64, limit int, offset int) ([]model.NotificationResponse, error) {
	notifications, err := s.notificationRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return s.convertToResponse(notifications), nil
}

func (s *notificationService) GetUnreadNotifications(userID int64) ([]model.NotificationResponse, error) {
	notifications, err := s.notificationRepo.FindUnreadByUserID(userID)
	if err != nil {
		return nil, err
	}

	return s.convertToResponse(notifications), nil
}

func (s *notificationService) GetUnreadCount(userID int64) (int64, error) {
	return s.notificationRepo.CountUnreadByUserID(userID)
}

func (s *notificationService) MarkAsRead(notificationIDs []int64, userID int64) error {
	return s.notificationRepo.MarkAsRead(notificationIDs, userID)
}

func (s *notificationService) MarkAllAsRead(userID int64) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}

func (s *notificationService) DeleteNotification(id int64, userID int64) error {
	notification, err := s.notificationRepo.FindByID(id)
	if err != nil {
		return err
	}

	if notification.UserID != userID {
		return errors.New("unauthorized to delete this notification")
	}

	return s.notificationRepo.Delete(id)
}

func (s *notificationService) CheckAndSendExpiryNotifications() error {
	now := time.Now()
	
	// Check untuk 30 hari sebelum expiry
	thirtyDaysLater := now.AddDate(0, 0, 30)
	// Check untuk 7 hari sebelum expiry
	sevenDaysLater := now.AddDate(0, 0, 7)
	// Check untuk hari ini (expired)
	today := now

	// Get permits yang akan expired dalam 30 hari
	permits30Days, err := s.permitRepo.FindExpiringPermits(now, thirtyDaysLater)
	if err != nil {
		return err
	}

	// Get permits yang akan expired dalam 7 hari
	permits7Days, err := s.permitRepo.FindExpiringPermits(now, sevenDaysLater)
	if err != nil {
		return err
	}

	// Get permits yang expired hari ini
	permitsExpired, err := s.permitRepo.FindExpiringPermits(today.AddDate(0, 0, -1), today)
	if err != nil {
		return err
	}

	// Process 30 days notification
	for _, permit := range permits30Days {
		daysLeft := int(time.Until(permit.ExpiryDate).Hours() / 24)
		if daysLeft <= 30 && daysLeft > 7 {
			s.processPermitNotification(permit, "expiry_reminder", daysLeft)
		}
	}

	// Process 7 days notification (warning)
	for _, permit := range permits7Days {
		daysLeft := int(time.Until(permit.ExpiryDate).Hours() / 24)
		if daysLeft <= 7 && daysLeft > 0 {
			s.processPermitNotification(permit, "expiry_warning", daysLeft)
		}
	}

	// Process expired notification
	for _, permit := range permitsExpired {
		s.processPermitNotification(permit, "expired", 0)
	}

	return nil
}

func (s *notificationService) processPermitNotification(permit model.Permit, notificationType string, daysLeft int) error {
	// Check apakah notifikasi sudah pernah dikirim dalam 24 jam terakhir
	yesterday := time.Now().AddDate(0, 0, -1)
	exists, err := s.notificationRepo.CheckExistingNotification(permit.ID, notificationType, yesterday)
	if err != nil {
		return err
	}
	if exists {
		return nil // Skip jika sudah ada notifikasi
	}

	// Collect recipients
	recipients := make(map[int64]*model.User)
	var emailRecipients []string

	// Add Responsible Person
	if permit.ResponsiblePersonID != nil && *permit.ResponsiblePersonID > 0 {
		user, err := s.userRepo.FindById(*permit.ResponsiblePersonID)
		if err == nil && user != nil {
			recipients[user.ID] = user
			emailRecipients = append(emailRecipients, user.Email)
		}
	}

	// Add Document Responsible Person
	if permit.ResponsibleDocPersonID != nil && *permit.ResponsibleDocPersonID > 0 {
		user, err := s.userRepo.FindById(*permit.ResponsibleDocPersonID)
		if err == nil && user != nil {
			if _, exists := recipients[user.ID]; !exists {
				recipients[user.ID] = user
				emailRecipients = append(emailRecipients, user.Email)
			}
		}
	}

	// Add Admin users (using role code)
	admins, err := s.userRepo.FindByRoleCode("ADMIN")
	if err == nil {
		for _, admin := range admins {
			if _, exists := recipients[admin.ID]; !exists {
				recipients[admin.ID] = &admin
				emailRecipients = append(emailRecipients, admin.Email)
			}
		}
	}

	// Add Permit Manager users from the same domain
	if permit.DomainID > 0 {
		managers, err := s.userRepo.FindByDomainAndRoleCode(permit.DomainID, "PERMIT_MANAGER")
		if err == nil {
			for _, manager := range managers {
				if _, exists := recipients[manager.ID]; !exists {
					recipients[manager.ID] = &manager
					emailRecipients = append(emailRecipients, manager.Email)
				}
			}
		}
	}

	// Create notification title and message
	title, message := s.getNotificationContent(permit, notificationType, daysLeft)

	// Create in-app notifications for each recipient
	for _, user := range recipients {
		notification := &model.Notification{
			UserID:   user.ID,
			PermitID: permit.ID,
			Type:     notificationType,
			Title:    title,
			Message:  message,
			IsRead:   false,
		}
		s.notificationRepo.Create(notification)
	}

	// Send email notification
	if len(emailRecipients) > 0 {
		expiryDateStr := permit.ExpiryDate.Format("02 January 2006")
		go helper.SendPermitExpiryNotification(
			emailRecipients,
			permit.Name,
			permit.PermitNo,
			expiryDateStr,
			daysLeft,
		)
	}

	return nil
}

func (s *notificationService) getNotificationContent(permit model.Permit, notificationType string, daysLeft int) (string, string) {
	switch notificationType {
	case "expiry_reminder":
		return fmt.Sprintf("Permit %s akan expired dalam %d hari", permit.Name, daysLeft),
			fmt.Sprintf("Permit %s (No: %s) akan expired pada %s. Segera lakukan perpanjangan.",
				permit.Name, permit.PermitNo, permit.ExpiryDate.Format("02 Jan 2006"))
	case "expiry_warning":
		return fmt.Sprintf("PERINGATAN: Permit %s akan expired dalam %d hari", permit.Name, daysLeft),
			fmt.Sprintf("Permit %s (No: %s) akan expired pada %s. Harap segera ditindaklanjuti!",
				permit.Name, permit.PermitNo, permit.ExpiryDate.Format("02 Jan 2006"))
	case "expired":
		return fmt.Sprintf("EXPIRED: Permit %s telah expired", permit.Name),
			fmt.Sprintf("Permit %s (No: %s) telah expired pada %s. Segera lakukan perpanjangan!",
				permit.Name, permit.PermitNo, permit.ExpiryDate.Format("02 Jan 2006"))
	default:
		return "Notifikasi Permit", "Ada update terkait permit Anda"
	}
}

func (s *notificationService) convertToResponse(notifications []model.Notification) []model.NotificationResponse {
	responses := make([]model.NotificationResponse, len(notifications))
	for i, notif := range notifications {
		responses[i] = model.NotificationResponse{
			ID:        notif.ID,
			UserID:    notif.UserID,
			PermitID:  notif.PermitID,
			Type:      notif.Type,
			Title:     notif.Title,
			Message:   notif.Message,
			IsRead:    notif.IsRead,
			ReadAt:    notif.ReadAt,
			CreatedAt: notif.CreatedAt,
		}

		if notif.Permit != nil {
			responses[i].Permit = &struct {
				ID         int64     `json:"id"`
				Name       string    `json:"name"`
				PermitNo   string    `json:"permit_no"`
				ExpiryDate time.Time `json:"expiry_date"`
			}{
				ID:         notif.Permit.ID,
				Name:       notif.Permit.Name,
				PermitNo:   notif.Permit.PermitNo,
				ExpiryDate: notif.Permit.ExpiryDate,
			}
		}
	}
	return responses
}
