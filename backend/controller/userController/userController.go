package userController

import (
	"net/http"
	"permit-app/helper/apiresponse"
	"permit-app/model"
	"permit-app/service/userService"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	service userService.UserService
}

func NewUserController(service userService.UserService) *UserController {
	return &UserController{service: service}
}

// Login authenticates a user and returns a JWT token
func (c *UserController) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	response, err := c.service.Login(&req)
	if err != nil {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid credentials", err, nil)
		return
	}

	apiresponse.OK(ctx, response, "Login successful", nil)
}

// SwitchDomain switches the user's current domain context and returns a new token
func (c *UserController) SwitchDomain(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in token", nil, nil)
		return
	}

	var req model.SwitchDomainRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	response, err := c.service.SwitchDomain(userID.(int64), &req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to switch domain", err, nil)
		return
	}

	apiresponse.OK(ctx, response, "Domain switched successfully", nil)
}

// GetProfile returns the profile of the currently logged-in user
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in token", nil, nil)
		return
	}

	user, err := c.service.GetUserByID(userID.(int64))
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve user profile", err, nil)
		return
	}

	apiresponse.OK(ctx, user, "Profile retrieved successfully", nil)
}

// UpdateProfile updates the profile of the currently logged-in user
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		apiresponse.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in token", nil, nil)
		return
	}

	var req model.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	user, err := c.service.UpdateProfile(userID.(int64), &req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to update profile", err, nil)
		return
	}

	apiresponse.OK(ctx, user, "Profile updated successfully", nil)
}

// Register creates a new user account
func (c *UserController) Register(ctx *gin.Context) {
	var req model.UserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	user, err := c.service.Register(&req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to register user", err, nil)
		return
	}

	apiresponse.Created(ctx, user, "User registered successfully", nil)
}

// GetByID retrieves a user by ID
func (c *UserController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	user, err := c.service.GetUserByID(id)
	if err != nil {
		apiresponse.Error(ctx, http.StatusNotFound, "NOT_FOUND", "User not found", err, nil)
		return
	}

	apiresponse.OK(ctx, user, "User retrieved successfully", nil)
}

// GetAll retrieves all users with filtering and pagination
func (c *UserController) GetAll(ctx *gin.Context) {
	var filter model.UserListRequest
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid query parameters", err, nil)
		return
	}

	if err := validator.New().Struct(&filter); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	// Extract domain_id from JWT token context
	if domainID, exists := ctx.Get("domain_id"); exists && domainID != nil {
		did := domainID.(int64)
		filter.DomainID = &did
	}

	users, total, err := c.service.GetAllUsers(&filter)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to retrieve users", err, nil)
		return
	}

	meta := apiresponse.PageMeta{
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: total,
	}
	apiresponse.OK(ctx, users, "Users retrieved successfully", meta)
}

// Update modifies an existing user
func (c *UserController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	var req model.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	user, err := c.service.UpdateUser(id, &req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to update user", err, nil)
		return
	}

	apiresponse.OK(ctx, user, "User updated successfully", nil)
}

// ChangePassword updates user password
func (c *UserController) ChangePassword(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	var req model.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	err = c.service.ChangePassword(id, &req)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to change password", err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Password changed successfully", nil)
}

// Delete removes a user
func (c *UserController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	err = c.service.DeleteUser(id)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to delete user", err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "User deleted successfully", nil)
}

// AddDomainRole adds a domain and role to a user
func (c *UserController) AddDomainRole(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	var req model.UserDomainRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid request body", err, nil)
		return
	}

	if err := validator.New().Struct(&req); err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Validation failed", err, nil)
		return
	}

	err = c.service.AddUserDomainRole(id, req.DomainID, req.RoleID, req.IsDefault)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to add domain-role to user", err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Domain-role added to user successfully", nil)
}

// RemoveDomainRole removes a domain and role from a user
func (c *UserController) RemoveDomainRole(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	domainID, err := strconv.ParseInt(ctx.Param("domain_id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid domain ID", err, nil)
		return
	}

	roleID, err := strconv.ParseInt(ctx.Param("role_id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid role ID", err, nil)
		return
	}

	err = c.service.RemoveUserDomainRole(id, domainID, roleID)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to remove domain-role from user", err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Domain-role removed from user successfully", nil)
}

// SetDefaultDomainRole sets the default domain-role for a user
func (c *UserController) SetDefaultDomainRole(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid ID", err, nil)
		return
	}

	domainID, err := strconv.ParseInt(ctx.Param("domain_id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid domain ID", err, nil)
		return
	}

	roleID, err := strconv.ParseInt(ctx.Param("role_id"), 10, 64)
	if err != nil {
		apiresponse.BadRequest(ctx, apiresponse.ErrCodeBadRequest, "Invalid role ID", err, nil)
		return
	}

	err = c.service.SetDefaultDomainRole(id, domainID, roleID)
	if err != nil {
		apiresponse.InternalServerError(ctx, apiresponse.ErrCodeInternal, "Failed to set default domain-role", err, nil)
		return
	}

	type EmptyData struct{}
	apiresponse.OK(ctx, EmptyData{}, "Default domain-role set successfully", nil)
}
