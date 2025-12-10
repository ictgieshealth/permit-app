package notificationController

import (
	"net/http"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/notificationService"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type NotificationController struct {
	notificationService notificationService.NotificationService
	validate            *validator.Validate
}

func NewNotificationController(notificationService notificationService.NotificationService, validate *validator.Validate) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
		validate:            validate,
	}
}

// GetNotifications godoc
// @Summary Get user notifications
// @Description Get paginated list of notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} apiresponse.Response
// @Failure 401 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Security BearerAuth
// @Router /notifications [get]
func (c *NotificationController) GetNotifications(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		apiresponse.Error(ctx, http.StatusUnauthorized, apiresponse.ErrCodeBadRequest, "Unauthorized", nil, nil)
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	notifications, err := c.notificationService.GetUserNotifications(userID.(int64), limit, offset)
	if err != nil {
		apiresponse.Error(ctx, http.StatusInternalServerError, apiresponse.ErrCodeInternal, "Failed to retrieve notifications", err, nil)
		return
	}

	apiresponse.OK(ctx, notifications, "Notifications retrieved successfully", nil)
}

// GetUnreadNotifications godoc
// @Summary Get unread notifications
// @Description Get all unread notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} apiresponse.Response
// @Failure 401 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Security BearerAuth
// @Router /notifications/unread [get]
func (c *NotificationController) GetUnreadNotifications(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		apiresponse.Error(ctx, http.StatusUnauthorized, apiresponse.ErrCodeBadRequest, "Unauthorized", nil, nil)
		return
	}

	notifications, err := c.notificationService.GetUnreadNotifications(userID.(int64))
	if err != nil {
		apiresponse.Error(ctx, http.StatusInternalServerError, apiresponse.ErrCodeInternal, "Failed to retrieve unread notifications", err, nil)
		return
	}

	apiresponse.OK(ctx, notifications, "Unread notifications retrieved successfully", nil)
}

// GetUnreadCount godoc
// @Summary Get unread notification count
// @Description Get count of unread notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} apiresponse.Response
// @Failure 401 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Security BearerAuth
// @Router /notifications/unread/count [get]
func (c *NotificationController) GetUnreadCount(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		apiresponse.Error(ctx, http.StatusUnauthorized, apiresponse.ErrCodeBadRequest, "Unauthorized", nil, nil)
		return
	}

	count, err := c.notificationService.GetUnreadCount(userID.(int64))
	if err != nil {
		apiresponse.Error(ctx, http.StatusInternalServerError, apiresponse.ErrCodeInternal, "Failed to retrieve unread count", err, nil)
		return
	}

	apiresponse.OK(ctx, gin.H{"count": count}, "Unread count retrieved successfully", nil)
}

// MarkAsRead godoc
// @Summary Mark notifications as read
// @Description Mark specific notifications as read
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body model.MarkAsReadRequest true "Notification IDs"
// @Success 200 {object} apiresponse.Response
// @Failure 400 {object} apiresponse.Response
// @Failure 401 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Security BearerAuth
// @Router /notifications/read [post]
func (c *NotificationController) MarkAsRead(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		apiresponse.Error(ctx, http.StatusUnauthorized, apiresponse.ErrCodeBadRequest, "Unauthorized", nil, nil)
		return
	}

	var request model.MarkAsReadRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		apiresponse.Error(ctx, http.StatusBadRequest, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := c.validate.Struct(request); err != nil {
		apiresponse.Error(ctx, http.StatusBadRequest, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	if err := c.notificationService.MarkAsRead(request.NotificationIDs, userID.(int64)); err != nil {
		apiresponse.Error(ctx, http.StatusInternalServerError, apiresponse.ErrCodeInternal, "Failed to mark notifications as read", err, nil)
		return
	}

	data := gin.H{"success": true}
	apiresponse.OK(ctx, data, "Notifications marked as read", nil)
}

// MarkAllAsRead godoc
// @Summary Mark all notifications as read
// @Description Mark all notifications as read for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} apiresponse.Response
// @Failure 401 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Security BearerAuth
// @Router /notifications/read/all [post]
func (c *NotificationController) MarkAllAsRead(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		apiresponse.Error(ctx, http.StatusUnauthorized, apiresponse.ErrCodeBadRequest, "Unauthorized", nil, nil)
		return
	}

	if err := c.notificationService.MarkAllAsRead(userID.(int64)); err != nil {
		apiresponse.Error(ctx, http.StatusInternalServerError, apiresponse.ErrCodeInternal, "Failed to mark all notifications as read", err, nil)
		return
	}

	data := gin.H{"success": true}
	apiresponse.OK(ctx, data, "All notifications marked as read", nil)
}

// DeleteNotification godoc
// @Summary Delete a notification
// @Description Delete a specific notification
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} apiresponse.Response
// @Failure 400 {object} apiresponse.Response
// @Failure 401 {object} apiresponse.Response
// @Failure 500 {object} apiresponse.Response
// @Security BearerAuth
// @Router /notifications/{id} [delete]
func (c *NotificationController) DeleteNotification(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		apiresponse.Error(ctx, http.StatusUnauthorized, apiresponse.ErrCodeBadRequest, "Unauthorized", nil, nil)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.Error(ctx, http.StatusBadRequest, apiresponse.ErrCodeBadRequest, "Invalid notification ID", err, nil)
		return
	}

	if err := c.notificationService.DeleteNotification(id, userID.(int64)); err != nil {
		apiresponse.Error(ctx, http.StatusInternalServerError, apiresponse.ErrCodeInternal, "Failed to delete notification", err, nil)
		return
	}

	data := gin.H{"success": true}
	apiresponse.OK(ctx, data, "Notification deleted successfully", nil)
}
