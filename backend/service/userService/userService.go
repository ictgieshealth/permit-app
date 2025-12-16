package userService

import (
	"errors"
	"permit-app/helper"
	"permit-app/model"
	"permit-app/repo/userRepository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	Register(req *model.UserRequest) (*model.UserResponse, error)
	Login(req *model.LoginRequest) (*model.LoginResponse, error)
	SwitchDomain(userID int64, req *model.SwitchDomainRequest) (*model.SwitchDomainResponse, error)
	GetUserByID(id int64) (*model.UserResponse, error)
	GetAllUsers(filter *model.UserListRequest) ([]model.UserResponse, int64, error)
	UpdateUser(id int64, req *model.UserUpdateRequest) (*model.UserResponse, error)
	UpdateProfile(id int64, req *model.UpdateProfileRequest) (*model.UserResponse, error)
	ChangePassword(id int64, req *model.ChangePasswordRequest) error
	DeleteUser(id int64) error
	AddUserDomainRole(userID int64, domainID int64, roleID int64, isDefault bool) error
	RemoveUserDomainRole(userID int64, domainID int64, roleID int64) error
	SetDefaultDomainRole(userID int64, domainID int64, roleID int64) error
}

type userService struct {
	repo userRepository.UserRepository
}

func NewUserService(repo userRepository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(req *model.UserRequest) (*model.UserResponse, error) {
	// Check if username already exists
	existingUsername, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUsername != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingEmail, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	user := &model.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashedPassword,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Nip:         req.Nip,
		IsActive:    isActive,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	// Create user-domain-role relationships
	for i, domainRole := range req.DomainRoles {
		userDomainRole := &model.UserDomainRole{
			UserID:    user.ID,
			DomainID:  domainRole.DomainID,
			RoleID:    domainRole.RoleID,
			IsDefault: i == 0 || domainRole.IsDefault, // First domain-role is default or as specified
		}
		err = s.repo.CreateUserDomainRole(userDomainRole)
		if err != nil {
			return nil, err
		}
	}

	// Reload user with associations
	created, err := s.repo.FindByID(user.ID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(created), nil
}

func (s *userService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	// Find user by username or email
	user, err := s.repo.FindByUsernameOrEmail(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Check password
	if !helper.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Get default domain-role or specified domain
	var selectedDomainRole *model.UserDomainRole
	if req.DomainID != nil {
		// Check if user has access to the specified domain
		userDomainRoles, _ := s.repo.GetUserDomainRoles(user.ID)
		for _, udr := range userDomainRoles {
			if udr.DomainID == *req.DomainID {
				selectedDomainRole = &udr
				break
			}
		}
		if selectedDomainRole == nil {
			return nil, errors.New("user does not have access to this domain")
		}
	} else {
		// Get default domain-role
		selectedDomainRole, err = s.repo.GetDefaultDomainRole(user.ID)
		if err != nil {
			return nil, err
		}
		if selectedDomainRole == nil {
			return nil, errors.New("user has no default domain-role set")
		}
	}

	// Generate JWT token with domain and role context
	token, err := helper.GenerateTokenWithDomain(user.ID, user.Username, user.Email, selectedDomainRole.DomainID, selectedDomainRole.RoleID)
	if err != nil {
		return nil, err
	}

	// Get all user domain-roles for domain switcher
	allDomainRoles, err := s.repo.GetUserDomainRoles(user.ID)
	if err != nil {
		return nil, err
	}

	var domainsResponse []model.UserDomainRoleResponse
	for _, udr := range allDomainRoles {
		domainsResponse = append(domainsResponse, *s.toUserDomainRoleResponse(&udr))
	}

	// Convert to response
	var currentDomain *model.DomainResponse
	var currentRole *model.RoleResponse
	if selectedDomainRole.Domain != nil {
		currentDomain = &model.DomainResponse{
			ID:          selectedDomainRole.Domain.ID,
			Code:        selectedDomainRole.Domain.Code,
			Name:        selectedDomainRole.Domain.Name,
			Description: selectedDomainRole.Domain.Description,
			IsActive:    selectedDomainRole.Domain.IsActive,
			CreatedAt:   selectedDomainRole.Domain.CreatedAt,
			UpdatedAt:   selectedDomainRole.Domain.UpdatedAt,
		}
	}
	if selectedDomainRole.Role != nil {
		currentRole = &model.RoleResponse{
			ID:          selectedDomainRole.Role.ID,
			Code:        selectedDomainRole.Role.Code,
			Name:        selectedDomainRole.Role.Name,
			Category:    selectedDomainRole.Role.Category,
			Description: selectedDomainRole.Role.Description,
			CreatedAt:   selectedDomainRole.Role.CreatedAt,
			UpdatedAt:   selectedDomainRole.Role.UpdatedAt,
		}
	}

	return &model.LoginResponse{
		Token:         token,
		User:          s.toResponse(user),
		CurrentDomain: currentDomain,
		CurrentRole:   currentRole,
		Domains:       domainsResponse,
	}, nil
}

// SwitchDomain switches user's current domain context
func (s *userService) SwitchDomain(userID int64, req *model.SwitchDomainRequest) (*model.SwitchDomainResponse, error) {
	// Get user
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Check if user has access to the requested domain
	userDomainRoles, err := s.repo.GetUserDomainRoles(userID)
	if err != nil {
		return nil, err
	}

	var selectedDomainRole *model.UserDomainRole
	for _, udr := range userDomainRoles {
		if udr.DomainID == req.DomainID {
			selectedDomainRole = &udr
			break
		}
	}

	if selectedDomainRole == nil {
		return nil, errors.New("user does not have access to this domain")
	}

	// Generate new JWT token with new domain context
	token, err := helper.GenerateTokenWithDomain(user.ID, user.Username, user.Email, selectedDomainRole.DomainID, selectedDomainRole.RoleID)
	if err != nil {
		return nil, err
	}

	// Build response
	var currentDomain *model.DomainResponse
	var currentRole *model.RoleResponse
	if selectedDomainRole.Domain != nil {
		currentDomain = &model.DomainResponse{
			ID:          selectedDomainRole.Domain.ID,
			Code:        selectedDomainRole.Domain.Code,
			Name:        selectedDomainRole.Domain.Name,
			Description: selectedDomainRole.Domain.Description,
			IsActive:    selectedDomainRole.Domain.IsActive,
			CreatedAt:   selectedDomainRole.Domain.CreatedAt,
			UpdatedAt:   selectedDomainRole.Domain.UpdatedAt,
		}
	}
	if selectedDomainRole.Role != nil {
		currentRole = &model.RoleResponse{
			ID:          selectedDomainRole.Role.ID,
			Code:        selectedDomainRole.Role.Code,
			Name:        selectedDomainRole.Role.Name,
			Category:    selectedDomainRole.Role.Category,
			Description:  selectedDomainRole.Role.Description,
			CreatedAt:   selectedDomainRole.Role.CreatedAt,
			UpdatedAt:   selectedDomainRole.Role.UpdatedAt,
		}
	}

	return &model.SwitchDomainResponse{
		Token:         token,
		CurrentDomain: currentDomain,
		CurrentRole:   currentRole,
	}, nil
}

func (s *userService) GetUserByID(id int64) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(user), nil
}

func (s *userService) GetAllUsers(filter *model.UserListRequest) ([]model.UserResponse, int64, error) {
	users, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.UserResponse
	for _, user := range users {
		responses = append(responses, *s.toResponse(&user))
	}

	return responses, total, nil
}

func (s *userService) UpdateUser(id int64, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check username uniqueness if changed
	if req.Username != "" && req.Username != user.Username {
		existing, err := s.repo.FindByUsername(req.Username)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("username already exists")
		}
		user.Username = req.Username
	}

	// Check email uniqueness if changed
	if req.Email != "" && req.Email != user.Email {
		existing, err := s.repo.FindByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := helper.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	if req.Nip != "" {
		user.Nip = req.Nip
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Update domain-role relationships if provided
	if req.DomainRoles != nil && len(req.DomainRoles) > 0 {
		// Delete existing domain-role relationships
		err = s.repo.DeleteUserDomainRoles(id)
		if err != nil {
			return nil, err
		}

		// Create new domain-role relationships
		for i, domainRole := range req.DomainRoles {
			userDomainRole := &model.UserDomainRole{
				UserID:    user.ID,
				DomainID:  domainRole.DomainID,
				RoleID:    domainRole.RoleID,
				IsDefault: i == 0 || domainRole.IsDefault, // First domain-role is default or as specified
			}
			err = s.repo.CreateUserDomainRole(userDomainRole)
			if err != nil {
				return nil, err
			}
		}
	}

	err = s.repo.Update(id, user)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updated), nil
}

func (s *userService) UpdateProfile(id int64, req *model.UpdateProfileRequest) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check email uniqueness if changed
	if req.Email != "" && req.Email != user.Email {
		existing, err := s.repo.FindByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	if req.Nip != "" {
		user.Nip = req.Nip
	}

	err = s.repo.Update(id, user)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updated), nil
}

func (s *userService) ChangePassword(id int64, req *model.ChangePasswordRequest) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Check old password
	if !helper.CheckPasswordHash(req.OldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.repo.Update(id, user)
}

func (s *userService) DeleteUser(id int64) error {
	return s.repo.Delete(id)
}

func (s *userService) AddUserDomainRole(userID int64, domainID int64, roleID int64, isDefault bool) error {
	userDomainRole := &model.UserDomainRole{
		UserID:    userID,
		DomainID:  domainID,
		RoleID:    roleID,
		IsDefault: isDefault,
	}
	return s.repo.CreateUserDomainRole(userDomainRole)
}

func (s *userService) RemoveUserDomainRole(userID int64, domainID int64, roleID int64) error {
	// Get user domain-roles
	userDomainRoles, err := s.repo.GetUserDomainRoles(userID)
	if err != nil {
		return err
	}

	// Check if this is the last domain-role
	if len(userDomainRoles) <= 1 {
		return errors.New("cannot remove the last domain-role from user")
	}

	// Delete the specific domain-role
	return s.repo.DeleteUserDomainRoles(userID)
}

func (s *userService) SetDefaultDomainRole(userID int64, domainID int64, roleID int64) error {
	return s.repo.UpdateDefaultDomainRole(userID, domainID, roleID)
}

func (s *userService) generateToken(user *model.User, domainID int64, roleID int64) (string, error) {
	jwtSecret := helper.GetEnv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"domain_id": domainID,
		"role_id":   roleID,
		"username":  user.Username,
		"email":     user.Email,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (s *userService) toResponse(user *model.User) *model.UserResponse {
	if user == nil {
		return nil
	}

	response := &model.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Nip:         user.Nip,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	// Convert user domain roles to response
	if user.UserDomainRoles != nil && len(user.UserDomainRoles) > 0 {
		for _, udr := range user.UserDomainRoles {
			response.DomainRoles = append(response.DomainRoles, *s.toUserDomainRoleResponse(&udr))
		}
	}

	return response
}

func (s *userService) toUserDomainRoleResponse(udr *model.UserDomainRole) *model.UserDomainRoleResponse {
	if udr == nil {
		return nil
	}

	response := &model.UserDomainRoleResponse{
		ID:        udr.ID,
		UserID:    udr.UserID,
		DomainID:  udr.DomainID,
		RoleID:    udr.RoleID,
		IsDefault: udr.IsDefault,
	}

	if udr.Domain != nil {
		response.Domain = &model.DomainResponse{
			ID:          udr.Domain.ID,
			Code:        udr.Domain.Code,
			Name:        udr.Domain.Name,
			Description: udr.Domain.Description,
			IsActive:    udr.Domain.IsActive,
			CreatedAt:   udr.Domain.CreatedAt,
			UpdatedAt:   udr.Domain.UpdatedAt,
		}
	}

	if udr.Role != nil {
		response.Role = &model.RoleResponse{
			ID:          udr.Role.ID,
			Code:        udr.Role.Code,
			Name:        udr.Role.Name,
			Category:    udr.Role.Category,
			Description: udr.Role.Description,
			CreatedAt:   udr.Role.CreatedAt,
			UpdatedAt:   udr.Role.UpdatedAt,
		}
	}

	return response
}
