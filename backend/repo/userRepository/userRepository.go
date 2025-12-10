package userRepository

import (
	"errors"
	"permit-app/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	CreateUserDomain(userDomain *model.UserDomain) error
	FindByID(id int64) (*model.User, error)
	FindById(id int64) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByUsernameOrEmail(usernameOrEmail string) (*model.User, error)
	FindAll(filter *model.UserListRequest) ([]model.User, int64, error)
	Update(id int64, user *model.User) error
	Delete(id int64) error
	DeleteUserDomains(userID int64) error
	GetUserDomains(userID int64) ([]model.UserDomain, error)
	GetDefaultDomain(userID int64) (*model.UserDomain, error)
	UpdateDefaultDomain(userID int64, domainID int64) error
	FindByRoleName(roleName string) ([]model.User, error)
	FindManagersByDomain(domainID int64) ([]model.User, error)
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

func (r *userRepository) CreateUserDomain(userDomain *model.UserDomain) error {
	return r.db.Create(userDomain).Error
}

func (r *userRepository) FindByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Role").
		Preload("Domains").
		Preload("UserDomains.Domain").
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
	err := r.db.Preload("Role").
		Preload("Domains").
		Preload("UserDomains.Domain").
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
	err := r.db.Preload("Role").
		Preload("Domains").
		Preload("UserDomains.Domain").
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
	err := r.db.Preload("Role").
		Preload("Domains").
		Preload("UserDomains.Domain").
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
		Preload("Role").
		Preload("Domains").
		Preload("UserDomains.Domain")

	// Filter by domain if specified
	if filter.DomainID != nil {
		query = query.Joins("JOIN user_domains ON user_domains.user_id = users.id").
			Where("user_domains.domain_id = ?", *filter.DomainID)
	}

	if filter.RoleID != nil {
		query = query.Where("role_id = ?", *filter.RoleID)
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

func (r *userRepository) DeleteUserDomains(userID int64) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.UserDomain{}).Error
}

func (r *userRepository) GetUserDomains(userID int64) ([]model.UserDomain, error) {
	var userDomains []model.UserDomain
	err := r.db.Preload("Domain").Where("user_id = ?", userID).Find(&userDomains).Error
	return userDomains, err
}

func (r *userRepository) GetDefaultDomain(userID int64) (*model.UserDomain, error) {
	var userDomain model.UserDomain
	err := r.db.Preload("Domain").
		Where("user_id = ? AND is_default = ?", userID, true).
		First(&userDomain).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userDomain, nil
}

func (r *userRepository) UpdateDefaultDomain(userID int64, domainID int64) error {
	// First, set all domains to non-default
	err := r.db.Model(&model.UserDomain{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
	if err != nil {
		return err
	}

	// Then, set the specified domain as default
	return r.db.Model(&model.UserDomain{}).
		Where("user_id = ? AND domain_id = ?", userID, domainID).
		Update("is_default", true).Error
}

func (r *userRepository) FindById(id int64) (*model.User, error) {
	return r.FindByID(id)
}

func (r *userRepository) FindByRoleName(roleName string) ([]model.User, error) {
	var users []model.User
	err := r.db.Joins("JOIN roles ON roles.id = users.role_id").
		Where("roles.name = ? AND users.is_active = ?", roleName, true).
		Preload("Role").
		Find(&users).Error
	return users, err
}

func (r *userRepository) FindManagersByDomain(domainID int64) ([]model.User, error) {
	var users []model.User
	err := r.db.Joins("JOIN roles ON roles.id = users.role_id").
		Joins("JOIN user_domains ON user_domains.user_id = users.id").
		Where("roles.name = ? AND user_domains.domain_id = ? AND users.is_active = ?", "manager", domainID, true).
		Preload("Role").
		Distinct().
		Find(&users).Error
	return users, err
}
