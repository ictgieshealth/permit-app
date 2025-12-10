package divisionRepository

import (
	"errors"
	"permit-app/model"

	"gorm.io/gorm"
)

type DivisionRepository interface {
	Create(division *model.Division) error
	FindByID(id int64) (*model.Division, error)
	FindByCodeAndDomainID(code string, domainID int64) (*model.Division, error)
	FindAll(filter *model.DivisionListRequest) ([]model.Division, int64, error)
	Update(id int64, division *model.Division) error
	Delete(id int64) error
}

type divisionRepository struct {
	db *gorm.DB
}

func NewDivisionRepository(db *gorm.DB) DivisionRepository {
	return &divisionRepository{db: db}
}

func (r *divisionRepository) Create(division *model.Division) error {
	return r.db.Create(division).Error
}

func (r *divisionRepository) FindByID(id int64) (*model.Division, error) {
	var division model.Division
	err := r.db.Preload("Domain").Where("id = ?", id).First(&division).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("division not found")
		}
		return nil, err
	}
	return &division, nil
}

func (r *divisionRepository) FindByCodeAndDomainID(code string, domainID int64) (*model.Division, error) {
	var division model.Division
	err := r.db.Where("code = ? AND domain_id = ?", code, domainID).First(&division).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &division, nil
}

func (r *divisionRepository) FindAll(filter *model.DivisionListRequest) ([]model.Division, int64, error) {
	var divisions []model.Division
	var total int64

	query := r.db.Model(&model.Division{}).Preload("Domain")

	if filter.DomainID != nil {
		query = query.Where("domain_id = ?", *filter.DomainID)
	}

	if filter.Code != "" {
		query = query.Where("code LIKE ?", "%"+filter.Code+"%")
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

	err = query.Order("created_at DESC").Find(&divisions).Error
	if err != nil {
		return nil, 0, err
	}

	return divisions, total, nil
}

func (r *divisionRepository) Update(id int64, division *model.Division) error {
	return r.db.Model(&model.Division{}).Where("id = ?", id).Updates(division).Error
}

func (r *divisionRepository) Delete(id int64) error {
	return r.db.Delete(&model.Division{}, id).Error
}