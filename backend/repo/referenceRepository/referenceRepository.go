package referenceRepository

import (
	"permit-app/model"

	"gorm.io/gorm"
)

type ReferenceRepository interface {
	Create(reference *model.Reference) error
	FindAll(filters map[string]interface{}, page, limit int) ([]model.Reference, int64, error)
	FindByID(id int64) (*model.Reference, error)
	FindByCategoryID(categoryID int64) ([]model.Reference, error)
	FindByModuleID(moduleID int64) ([]model.Reference, error)
	Update(reference *model.Reference) error
	Delete(id int64) error
}

type referenceRepositoryImpl struct {
	db *gorm.DB
}

func NewReferenceRepository(db *gorm.DB) ReferenceRepository {
	return &referenceRepositoryImpl{db: db}
}

func (r *referenceRepositoryImpl) Create(reference *model.Reference) error {
	return r.db.Create(reference).Error
}

func (r *referenceRepositoryImpl) FindAll(filters map[string]interface{}, page, limit int) ([]model.Reference, int64, error) {
	var references []model.Reference
	var total int64

	query := r.db.Model(&model.Reference{}).
		Preload("ReferenceCategory").
		Preload("ReferenceCategory.Module")

	// Apply filters
	if categoryID, ok := filters["reference_category_id"].(int64); ok && categoryID > 0 {
		query = query.Where("reference_category_id = ?", categoryID)
	}
	if moduleID, ok := filters["module_id"].(int64); ok && moduleID > 0 {
		query = query.Joins("JOIN reference_categories ON references.reference_category_id = reference_categories.id").
			Where("reference_categories.module_id = ?", moduleID)
	}
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("references.name ILIKE ?", "%"+name+"%")
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("references.is_active = ?", isActive)
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
	if err := query.Order("references.id ASC").Find(&references).Error; err != nil {
		return nil, 0, err
	}

	return references, total, nil
}

func (r *referenceRepositoryImpl) FindByID(id int64) (*model.Reference, error) {
	var reference model.Reference
	if err := r.db.Preload("ReferenceCategory").
		Preload("ReferenceCategory.Module").
		First(&reference, id).Error; err != nil {
		return nil, err
	}
	return &reference, nil
}

func (r *referenceRepositoryImpl) FindByCategoryID(categoryID int64) ([]model.Reference, error) {
	var references []model.Reference
	if err := r.db.Where("reference_category_id = ? AND is_active = ?", categoryID, true).
		Preload("ReferenceCategory").
		Preload("ReferenceCategory.Module").
		Order("id ASC").
		Find(&references).Error; err != nil {
		return nil, err
	}
	return references, nil
}

func (r *referenceRepositoryImpl) FindByModuleID(moduleID int64) ([]model.Reference, error) {
	var references []model.Reference
	if err := r.db.Joins("JOIN reference_categories ON references.reference_category_id = reference_categories.id").
		Where("reference_categories.module_id = ? AND references.is_active = ?", moduleID, true).
		Preload("ReferenceCategory").
		Preload("ReferenceCategory.Module").
		Order("references.id ASC").
		Find(&references).Error; err != nil {
		return nil, err
	}
	return references, nil
}

func (r *referenceRepositoryImpl) Update(reference *model.Reference) error {
	return r.db.Save(reference).Error
}

func (r *referenceRepositoryImpl) Delete(id int64) error {
	return r.db.Delete(&model.Reference{}, id).Error
}
