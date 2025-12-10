package permitRepository

import (
	"errors"
	"permit-app/model"
	"time"

	"gorm.io/gorm"
)

type PermitRepository interface {
	Create(permit *model.Permit) error
	FindByID(id int64) (*model.Permit, error)
	FindByPermitNoAndDomainID(permitNo string, domainID int64) (*model.Permit, error)
	FindAll(filter *model.PermitListRequest) ([]model.Permit, int64, error)
	Update(id int64, permit *model.Permit) error
	Delete(id int64) error
	FindExpiringPermits(startDate time.Time, endDate time.Time) ([]model.Permit, error)
	Search(query string, filter *model.PermitListRequest) ([]model.Permit, int64, error)
}

type permitRepository struct {
	db *gorm.DB
}

func NewPermitRepository(db *gorm.DB) PermitRepository {
	return &permitRepository{db: db}
}

func (r *permitRepository) Create(permit *model.Permit) error {
	return r.db.Create(permit).Error
}

func (r *permitRepository) FindByID(id int64) (*model.Permit, error) {
	var permit model.Permit
	err := r.db.Preload("Domain").Preload("Division.Domain").Preload("PermitType.Division.Domain").Preload("ResponsiblePerson").Preload("ResponsibleDocPerson").Where("id = ?", id).First(&permit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permit not found")
		}
		return nil, err
	}
	return &permit, nil
}

func (r *permitRepository) FindByPermitNoAndDomainID(permitNo string, domainID int64) (*model.Permit, error) {
	var permit model.Permit
	err := r.db.Where("permit_no = ? AND domain_id = ?", permitNo, domainID).First(&permit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permit, nil
}

func (r *permitRepository) FindAll(filter *model.PermitListRequest) ([]model.Permit, int64, error) {
	var permits []model.Permit
	var total int64

	query := r.db.Model(&model.Permit{}).Preload("Domain").Preload("Division.Domain").Preload("PermitType.Division.Domain").Preload("ResponsiblePerson").Preload("ResponsibleDocPerson")

	if filter.DomainID != nil {
		query = query.Where("domain_id = ?", *filter.DomainID)
	}

	if filter.DivisionID != nil {
		query = query.Where("division_id = ?", *filter.DivisionID)
	}

	if filter.PermitTypeID != nil {
		query = query.Where("permit_type_id = ?", *filter.PermitTypeID)
	}

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.ApplicationType != "" {
		query = query.Where("application_type = ?", filter.ApplicationType)
	}

	if filter.PermitNo != "" {
		query = query.Where("permit_no LIKE ?", "%"+filter.PermitNo+"%")
	}

	if filter.ResponsiblePerson != "" {
		query = query.Joins("LEFT JOIN users ON users.id = permits.responsible_person_id").Where("users.full_name LIKE ?", "%"+filter.ResponsiblePerson+"%")
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	err = query.Order("created_at DESC").Find(&permits).Error
	if err != nil {
		return nil, 0, err
	}

	return permits, total, nil
}

func (r *permitRepository) Update(id int64, permit *model.Permit) error {
	return r.db.Model(&model.Permit{}).Where("id = ?", id).Updates(permit).Error
}

func (r *permitRepository) Delete(id int64) error {
	return r.db.Delete(&model.Permit{}, id).Error
}

func (r *permitRepository) FindExpiringPermits(startDate time.Time, endDate time.Time) ([]model.Permit, error) {
var permits []model.Permit
err := r.db.Where("expiry_date BETWEEN ? AND ? AND status = ?", startDate, endDate, "active").
Preload("Domain").
Preload("ResponsiblePerson").
Preload("ResponsibleDocPerson").
Find(&permits).Error
return permits, err
}

func (r *permitRepository) Search(query string, filter *model.PermitListRequest) ([]model.Permit, int64, error) {
	var permits []model.Permit
	var total int64

	db := r.db.Model(&model.Permit{}).Preload("Domain").Preload("Division.Domain").Preload("PermitType.Division.Domain").Preload("ResponsiblePerson").Preload("ResponsibleDocPerson")

	// Search across multiple fields
	searchPattern := "%" + query + "%"
	db = db.Where("name LIKE ? OR permit_no LIKE ? OR doc_name LIKE ? OR doc_number LIKE ?", searchPattern, searchPattern, searchPattern, searchPattern)

	// Apply additional filters if provided
	if filter.DomainID != nil {
		db = db.Where("domain_id = ?", *filter.DomainID)
	}

	if filter.DivisionID != nil {
		db = db.Where("division_id = ?", *filter.DivisionID)
	}

	if filter.PermitTypeID != nil {
		db = db.Where("permit_type_id = ?", *filter.PermitTypeID)
	}

	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		db = db.Offset(offset).Limit(filter.Limit)
	}

	err = db.Order("created_at DESC").Find(&permits).Error
	if err != nil {
		return nil, 0, err
	}

	return permits, total, nil
}
