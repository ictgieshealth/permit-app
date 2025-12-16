package roleRepository

import (
	"errors"
	"permit-app/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *model.Role) error
	FindByID(id int64) (*model.Role, error)
	FindByCode(code string) (*model.Role, error)
	FindByName(name string) (*model.Role, error)
	FindAll(filter *model.RoleListRequest) ([]model.Role, int64, error)
	Update(id int64, role *model.Role) error
	Delete(id int64) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) FindByID(id int64) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByCode(code string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("code = ?", code).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindAll(filter *model.RoleListRequest) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	query := r.db.Model(&model.Role{})

	if filter.Code != "" {
		query = query.Where("code LIKE ?", "%"+filter.Code+"%")
	}

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	err = query.Order("category ASC, name ASC").Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *roleRepository) Update(id int64, role *model.Role) error {
	return r.db.Where("id = ?", id).Updates(role).Error
}

func (r *roleRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.Role{}).Error
}
