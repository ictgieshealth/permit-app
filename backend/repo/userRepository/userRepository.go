package userRepository

import (
	"errors"
	"permit-app/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	CreateUserDomainRole(userDomainRole *model.UserDomainRole) error
	FindByID(id int64) (*model.User, error)
	FindById(id int64) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByUsernameOrEmail(usernameOrEmail string) (*model.User, error)
	FindAll(filter *model.UserListRequest) ([]model.User, int64, error)
	Update(id int64, user *model.User) error
	Delete(id int64) error
	DeleteUserDomainRoles(userID int64) error
	GetUserDomainRoles(userID int64) ([]model.UserDomainRole, error)
	GetDefaultDomainRole(userID int64) (*model.UserDomainRole, error)
	UpdateDefaultDomainRole(userID int64, domainID int64, roleID int64) error
	FindByRoleCode(roleCode string) ([]model.User, error)
	FindByDomainAndRoleCode(domainID int64, roleCode string) ([]model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) CreateUserDomainRole(userDomainRole *model.UserDomainRole) error {
	return r.db.Create(userDomainRole).Error
}

func (r *userRepository) FindByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.Preload("UserDomainRoles.Domain").
		Preload("UserDomainRoles.Role").
		Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("UserDomainRoles.Domain").
		Preload("UserDomainRoles.Role").
		Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("UserDomainRoles.Domain").
		Preload("UserDomainRoles.Role").
		Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsernameOrEmail(usernameOrEmail string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("UserDomainRoles.Domain").
		Preload("UserDomainRoles.Role").
		Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(filter *model.UserListRequest) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{}).
		Preload("UserDomainRoles.Domain").
		Preload("UserDomainRoles.Role")

	// Filter by domain if specified
	if filter.DomainID != nil {
		query = query.Joins("JOIN user_domain_roles ON user_domain_roles.user_id = users.id").
			Where("user_domain_roles.domain_id = ?", *filter.DomainID)
	}

	// Filter by role if specified
	if filter.RoleID != nil {
		if filter.DomainID == nil {
			query = query.Joins("JOIN user_domain_roles ON user_domain_roles.user_id = users.id")
		}
		query = query.Where("user_domain_roles.role_id = ?", *filter.RoleID)
	}

	// Filter by role category if specified
	if filter.Category != "" {
		if filter.DomainID == nil && filter.RoleID == nil {
			query = query.Joins("JOIN user_domain_roles ON user_domain_roles.user_id = users.id")
		}
		query = query.Joins("JOIN roles ON roles.id = user_domain_roles.role_id").
			Where("roles.category = ?", filter.Category)
	}

	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}

	if filter.Email != "" {
		query = query.Where("email LIKE ?", "%"+filter.Email+"%")
	}

	if filter.FullName != "" {
		query = query.Where("full_name LIKE ?", "%"+filter.FullName+"%")
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// Use distinct to avoid duplicate users when joining
	query = query.Distinct()

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	err = query.Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(id int64, user *model.User) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(user).Error
}

func (r *userRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *userRepository) DeleteUserDomainRoles(userID int64) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.UserDomainRole{}).Error
}

func (r *userRepository) GetUserDomainRoles(userID int64) ([]model.UserDomainRole, error) {
	var userDomainRoles []model.UserDomainRole
	err := r.db.Preload("Domain").
		Preload("Role").
		Where("user_id = ?", userID).
		Find(&userDomainRoles).Error
	return userDomainRoles, err
}

func (r *userRepository) GetDefaultDomainRole(userID int64) (*model.UserDomainRole, error) {
	var userDomainRole model.UserDomainRole
	err := r.db.Preload("Domain").
		Preload("Role").
		Where("user_id = ? AND is_default = ?", userID, true).
		First(&userDomainRole).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userDomainRole, nil
}

func (r *userRepository) UpdateDefaultDomainRole(userID int64, domainID int64, roleID int64) error {
	// First, set all domain-roles to non-default for this user
	err := r.db.Model(&model.UserDomainRole{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
	if err != nil {
		return err
	}

	// Then, set the specified domain-role as default
	return r.db.Model(&model.UserDomainRole{}).
		Where("user_id = ? AND domain_id = ? AND role_id = ?", userID, domainID, roleID).
		Update("is_default", true).Error
}

func (r *userRepository) FindById(id int64) (*model.User, error) {
	return r.FindByID(id)
}

func (r *userRepository) FindByRoleCode(roleCode string) ([]model.User, error) {
	var users []model.User
	err := r.db.Joins("JOIN user_domain_roles ON user_domain_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_domain_roles.role_id").
		Where("roles.code = ? AND users.is_active = ?", roleCode, true).
		Preload("UserDomainRoles.Domain").
		Preload("UserDomainRoles.Role").
		Distinct().
		Find(&users).Error
	return users, err
}

func (r *userRepository) FindByDomainAndRoleCode(domainID int64, roleCode string) ([]model.User, error) {
	var users []model.User
	err := r.db.Joins("JOIN user_domain_roles ON user_domain_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_domain_roles.role_id").
		Where("roles.code = ? AND user_domain_roles.domain_id = ? AND users.is_active = ?", roleCode, domainID, true).
		Preload("UserDomainRoles.Domain").
		Preload("UserDomainRoles.Role").
		Distinct().
		Find(&users).Error
	return users, err
}
