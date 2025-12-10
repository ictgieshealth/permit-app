package domainRepository

import (
	"errors"
	"permit-app/model"

	"gorm.io/gorm"
)

type DomainRepository interface {
	Create(domain *model.Domain) error
	FindByID(id int64) (*model.Domain, error)
	FindByCode(code string) (*model.Domain, error)
	FindAll(filter *model.DomainListRequest) ([]model.Domain, int64, error)
	Update(id int64, domain *model.Domain) error
	Delete(id int64) error
}

type domainRepository struct {
	db *gorm.DB
}

func NewDomainRepository(db *gorm.DB) DomainRepository {
	return &domainRepository{db: db}
}

func (r *domainRepository) Create(domain *model.Domain) error {
	return r.db.Create(domain).Error
}

func (r *domainRepository) FindByID(id int64) (*model.Domain, error) {
	var domain model.Domain
	err := r.db.Where("id = ?", id).First(&domain).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("domain not found")
		}
		return nil, err
	}
	return &domain, nil
}

func (r *domainRepository) FindByCode(code string) (*model.Domain, error) {
	var domain model.Domain
	err := r.db.Where("code = ?", code).First(&domain).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &domain, nil
}

func (r *domainRepository) FindAll(filter *model.DomainListRequest) ([]model.Domain, int64, error) {
	var domains []model.Domain
	var total int64

	query := r.db.Model(&model.Domain{})

	if filter.Code != "" {
		query = query.Where("code LIKE ?", "%"+filter.Code+"%")
	}

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
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

	err = query.Order("created_at DESC").Find(&domains).Error
	if err != nil {
		return nil, 0, err
	}

	return domains, total, nil
}

func (r *domainRepository) Update(id int64, domain *model.Domain) error {
	return r.db.Where("id = ?", id).Updates(domain).Error
}

func (r *domainRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.Domain{}).Error
}
