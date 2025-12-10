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
	GetUserByID(id int64) (*model.UserResponse, error)
	GetAllUsers(filter *model.UserListRequest) ([]model.UserResponse, int64, error)
	UpdateUser(id int64, req *model.UserUpdateRequest) (*model.UserResponse, error)
	UpdateProfile(id int64, req *model.UpdateProfileRequest) (*model.UserResponse, error)
	ChangePassword(id int64, req *model.ChangePasswordRequest) error
	DeleteUser(id int64) error
	AddUserDomain(userID int64, domainID int64, isDefault bool) error
	RemoveUserDomain(userID int64, domainID int64) error
	SetDefaultDomain(userID int64, domainID int64) error
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
		RoleID:      req.RoleID,
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashedPassword,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		NIP:         req.NIP,
		IsActive:    isActive,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	// Create user-domain relationships
	for i, domainID := range req.DomainIDs {
		userDomain := &model.UserDomain{
			UserID:    user.ID,
			DomainID:  domainID,
			IsDefault: i == 0, // First domain is default
		}
		err = s.repo.CreateUserDomain(userDomain)
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

	// Get default domain or specified domain
	var selectedDomain *model.UserDomain
	if req.DomainID != nil {
		// Check if user has access to the specified domain
		userDomains, _ := s.repo.GetUserDomains(user.ID)
		for _, ud := range userDomains {
			if ud.DomainID == *req.DomainID {
				selectedDomain = &ud
				break
			}
		}
		if selectedDomain == nil {
			return nil, errors.New("user does not have access to this domain")
		}
	} else {
		// Get default domain
		selectedDomain, err = s.repo.GetDefaultDomain(user.ID)
		if err != nil {
			return nil, err
		}
		if selectedDomain == nil {
			return nil, errors.New("user has no default domain set")
		}
	}

	// Generate JWT token with domain context
	token, err := s.generateToken(user, selectedDomain.DomainID)
	if err != nil {
		return nil, err
	}

	// Convert domain to response
	var defaultDomainResponse *model.DomainResponse
	if selectedDomain.Domain != nil {
		defaultDomainResponse = &model.DomainResponse{
			ID:          selectedDomain.Domain.ID,
			Code:        selectedDomain.Domain.Code,
			Name:        selectedDomain.Domain.Name,
			Description: selectedDomain.Domain.Description,
			IsActive:    selectedDomain.Domain.IsActive,
			CreatedAt:   selectedDomain.Domain.CreatedAt,
			UpdatedAt:   selectedDomain.Domain.UpdatedAt,
		}
	}

	return &model.LoginResponse{
		Token:         token,
		User:          s.toResponse(user),
		DefaultDomain: defaultDomainResponse,
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

	if req.RoleID > 0 {
		user.RoleID = req.RoleID
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

	if req.NIP != "" {
		user.NIP = req.NIP
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Update domain relationships if provided
	if req.DomainIDs != nil && len(req.DomainIDs) > 0 {
		// Delete existing domain relationships
		err = s.repo.DeleteUserDomains(id)
		if err != nil {
			return nil, err
		}

		// Create new domain relationships
		for i, domainID := range req.DomainIDs {
			userDomain := &model.UserDomain{
				UserID:    user.ID,
				DomainID:  domainID,
				IsDefault: i == 0, // First domain is default
			}
			err = s.repo.CreateUserDomain(userDomain)
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

	if req.NIP != "" {
		user.NIP = req.NIP
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

func (s *userService) AddUserDomain(userID int64, domainID int64, isDefault bool) error {
	userDomain := &model.UserDomain{
		UserID:    userID,
		DomainID:  domainID,
		IsDefault: isDefault,
	}
	return s.repo.CreateUserDomain(userDomain)
}

func (s *userService) RemoveUserDomain(userID int64, domainID int64) error {
	// Get user domains
	userDomains, err := s.repo.GetUserDomains(userID)
	if err != nil {
		return err
	}

	// Check if this is the last domain
	if len(userDomains) <= 1 {
		return errors.New("cannot remove the last domain from user")
	}

	// Delete the specific domain
	return s.repo.DeleteUserDomains(userID)
}

func (s *userService) SetDefaultDomain(userID int64, domainID int64) error {
	return s.repo.UpdateDefaultDomain(userID, domainID)
}

func (s *userService) generateToken(user *model.User, domainID int64) (string, error) {
	jwtSecret := helper.GetEnv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"domain_id": domainID,
		"role_id":   user.RoleID,
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
		RoleID:      user.RoleID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		NIP:         user.NIP,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	if user.Role != nil {
		response.Role = &model.RoleResponse{
			ID:        user.Role.ID,
			Name:      user.Role.Name,
			CreatedAt: user.Role.CreatedAt,
			UpdatedAt: user.Role.UpdatedAt,
		}
	}

	// Convert domains to response
	if user.Domains != nil && len(user.Domains) > 0 {
		for _, domain := range user.Domains {
			response.Domains = append(response.Domains, model.DomainResponse{
				ID:          domain.ID,
				Code:        domain.Code,
				Name:        domain.Name,
				Description: domain.Description,
				IsActive:    domain.IsActive,
				CreatedAt:   domain.CreatedAt,
				UpdatedAt:   domain.UpdatedAt,
			})
		}
	}

	return response
}
