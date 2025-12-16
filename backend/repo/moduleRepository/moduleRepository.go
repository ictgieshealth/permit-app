package moduleRepository

import (
	"permit-app/model"

	"gorm.io/gorm"
)

type ModuleRepository interface {
	Create(module *model.Module) error
	FindAll(filters map[string]interface{}, page, limit int) ([]model.Module, int64, error)
	FindByID(id int64) (*model.Module, error)
	FindByCode(code string) (*model.Module, error)
	Update(module *model.Module) error
	Delete(id int64) error
}

type moduleRepositoryImpl struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) ModuleRepository {
	return &moduleRepositoryImpl{db: db}
}

func (r *moduleRepositoryImpl) Create(module *model.Module) error {
	return r.db.Create(module).Error
}

func (r *moduleRepositoryImpl) FindAll(filters map[string]interface{}, page, limit int) ([]model.Module, int64, error) {
	var modules []model.Module
	var total int64

	query := r.db.Model(&model.Module{})

	// Apply filters
	if code, ok := filters["code"].(string); ok && code != "" {
		query = query.Where("code ILIKE ?", "%"+code+"%")
	}
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	// Execute query
	if err := query.Order("id ASC").Find(&modules).Error; err != nil {
		return nil, 0, err
	}

	return modules, total, nil
}

func (r *moduleRepositoryImpl) FindByID(id int64) (*model.Module, error) {
	var module model.Module
	if err := r.db.First(&module, id).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

func (r *moduleRepositoryImpl) FindByCode(code string) (*model.Module, error) {
	var module model.Module
	if err := r.db.Where("code = ?", code).First(&module).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

func (r *moduleRepositoryImpl) Update(module *model.Module) error {
	return r.db.Save(module).Error
}

func (r *moduleRepositoryImpl) Delete(id int64) error {
	return r.db.Delete(&model.Module{}, id).Error
}
