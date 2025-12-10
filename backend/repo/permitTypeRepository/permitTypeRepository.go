package permitTypeRepository

import (
	"errors"
	"permit-app/model"

	"gorm.io/gorm"
)

type PermitTypeRepository interface {
	Create(permitType *model.PermitType) error
	FindByID(id int64) (*model.PermitType, error)
	FindAll(filter *model.PermitTypeListRequest) ([]model.PermitType, int64, error)
	Update(id int64, permitType *model.PermitType) error
	Delete(id int64) error
}

type permitTypeRepository struct {
	db *gorm.DB
}

func NewPermitTypeRepository(db *gorm.DB) PermitTypeRepository {
	return &permitTypeRepository{db: db}
}

func (r *permitTypeRepository) Create(permitType *model.PermitType) error {
	return r.db.Create(permitType).Error
}

func (r *permitTypeRepository) FindByID(id int64) (*model.PermitType, error) {
	var permitType model.PermitType
	err := r.db.Preload("Division.Domain").Where("id = ?", id).First(&permitType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permit type not found")
		}
		return nil, err
	}
	return &permitType, nil
}

func (r *permitTypeRepository) FindAll(filter *model.PermitTypeListRequest) ([]model.PermitType, int64, error) {
	var permitTypes []model.PermitType
	var total int64

	query := r.db.Model(&model.PermitType{}).Preload("Division.Domain")

	if filter.DivisionID != nil && *filter.DivisionID > 0 {
		query = query.Where("division_id = ?", *filter.DivisionID)
	}

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	err = query.Order("created_at DESC").Find(&permitTypes).Error
	if err != nil {
		return nil, 0, err
	}

	return permitTypes, total, nil
}

func (r *permitTypeRepository) Update(id int64, permitType *model.PermitType) error {
	return r.db.Model(&model.PermitType{}).Where("id = ?", id).Updates(permitType).Error
}

func (r *permitTypeRepository) Delete(id int64) error {
	return r.db.Delete(&model.PermitType{}, id).Error
}