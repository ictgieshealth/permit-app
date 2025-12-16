package referenceCategoryRepository

import (
	"permit-app/model"

	"gorm.io/gorm"
)

type ReferenceCategoryRepository interface {
	Create(category *model.ReferenceCategory) error
	FindAll(filters map[string]interface{}, page, limit int) ([]model.ReferenceCategory, int64, error)
	FindByID(id int64) (*model.ReferenceCategory, error)
	FindByModuleID(moduleID int64) ([]model.ReferenceCategory, error)
	Update(category *model.ReferenceCategory) error
	Delete(id int64) error
}

type referenceCategoryRepositoryImpl struct {
	db *gorm.DB
}

func NewReferenceCategoryRepository(db *gorm.DB) ReferenceCategoryRepository {
	return &referenceCategoryRepositoryImpl{db: db}
}

func (r *referenceCategoryRepositoryImpl) Create(category *model.ReferenceCategory) error {
	return r.db.Create(category).Error
}

func (r *referenceCategoryRepositoryImpl) FindAll(filters map[string]interface{}, page, limit int) ([]model.ReferenceCategory, int64, error) {
	var categories []model.ReferenceCategory
	var total int64

	query := r.db.Model(&model.ReferenceCategory{}).Preload("Module")

	// Apply filters
	if moduleID, ok := filters["module_id"].(int64); ok && moduleID > 0 {
		query = query.Where("module_id = ?", moduleID)
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
	if err := query.Order("id ASC").Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *referenceCategoryRepositoryImpl) FindByID(id int64) (*model.ReferenceCategory, error) {
	var category model.ReferenceCategory
	if err := r.db.Preload("Module").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *referenceCategoryRepositoryImpl) FindByModuleID(moduleID int64) ([]model.ReferenceCategory, error) {
	var categories []model.ReferenceCategory
	if err := r.db.Where("module_id = ? AND is_active = ?", moduleID, true).
		Preload("Module").
		Order("id ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *referenceCategoryRepositoryImpl) Update(category *model.ReferenceCategory) error {
	return r.db.Save(category).Error
}

func (r *referenceCategoryRepositoryImpl) Delete(id int64) error {
	return r.db.Delete(&model.ReferenceCategory{}, id).Error
}
