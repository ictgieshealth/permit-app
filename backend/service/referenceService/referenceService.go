package referenceService

import (
	"errors"
	"permit-app/model"
	"permit-app/repo/referenceCategoryRepository"
	"permit-app/repo/referenceRepository"

	"gorm.io/gorm"
)

type ReferenceService interface {
	CreateReference(req model.ReferenceRequest) (*model.ReferenceResponse, error)
	GetReferences(req model.ReferenceListRequest) ([]model.ReferenceResponse, int64, error)
	GetReferenceByID(id int64) (*model.ReferenceResponse, error)
	GetReferencesByCategoryID(categoryID int64) ([]model.ReferenceResponse, error)
	GetReferencesByModuleID(moduleID int64) ([]model.ReferenceResponse, error)
	UpdateReference(id int64, req model.ReferenceRequest) (*model.ReferenceResponse, error)
	DeleteReference(id int64) error
}

type referenceServiceImpl struct {
	referenceRepo referenceRepository.ReferenceRepository
	categoryRepo  referenceCategoryRepository.ReferenceCategoryRepository
}

func NewReferenceService(
	referenceRepo referenceRepository.ReferenceRepository,
	categoryRepo referenceCategoryRepository.ReferenceCategoryRepository,
) ReferenceService {
	return &referenceServiceImpl{
		referenceRepo: referenceRepo,
		categoryRepo:  categoryRepo,
	}
}

func (s *referenceServiceImpl) CreateReference(req model.ReferenceRequest) (*model.ReferenceResponse, error) {
	// Validate category exists
	_, err := s.categoryRepo.FindByID(req.ReferenceCategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reference category not found")
		}
		return nil, err
	}

	reference := &model.Reference{
		ReferenceCategoryID: req.ReferenceCategoryID,
		Name:                req.Name,
		IsActive:            true,
	}

	if req.IsActive != nil {
		reference.IsActive = *req.IsActive
	}

	if err := s.referenceRepo.Create(reference); err != nil {
		return nil, err
	}

	// Reload with relations
	reference, err = s.referenceRepo.FindByID(reference.ID)
	if err != nil {
		return nil, err
	}

	response := model.ToReferenceResponse(reference)
	return &response, nil
}

func (s *referenceServiceImpl) GetReferences(req model.ReferenceListRequest) ([]model.ReferenceResponse, int64, error) {
	filters := make(map[string]interface{})

	if req.ReferenceCategoryID > 0 {
		filters["reference_category_id"] = req.ReferenceCategoryID
	}
	if req.ModuleID > 0 {
		filters["module_id"] = req.ModuleID
	}
	if req.Name != "" {
		filters["name"] = req.Name
	}
	if req.IsActive != nil {
		filters["is_active"] = *req.IsActive
	}

	references, total, err := s.referenceRepo.FindAll(filters, req.Page, req.Limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response
	responses := make([]model.ReferenceResponse, len(references))
	for i, reference := range references {
		responses[i] = model.ToReferenceResponse(&reference)
	}

	return responses, total, nil
}

func (s *referenceServiceImpl) GetReferenceByID(id int64) (*model.ReferenceResponse, error) {
	reference, err := s.referenceRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reference not found")
		}
		return nil, err
	}

	response := model.ToReferenceResponse(reference)
	return &response, nil
}

func (s *referenceServiceImpl) GetReferencesByCategoryID(categoryID int64) ([]model.ReferenceResponse, error) {
	references, err := s.referenceRepo.FindByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ReferenceResponse, len(references))
	for i, reference := range references {
		responses[i] = model.ToReferenceResponse(&reference)
	}

	return responses, nil
}

func (s *referenceServiceImpl) GetReferencesByModuleID(moduleID int64) ([]model.ReferenceResponse, error) {
	references, err := s.referenceRepo.FindByModuleID(moduleID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ReferenceResponse, len(references))
	for i, reference := range references {
		responses[i] = model.ToReferenceResponse(&reference)
	}

	return responses, nil
}

func (s *referenceServiceImpl) UpdateReference(id int64, req model.ReferenceRequest) (*model.ReferenceResponse, error) {
	reference, err := s.referenceRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reference not found")
		}
		return nil, err
	}

	// Validate category exists if changed
	if req.ReferenceCategoryID != reference.ReferenceCategoryID {
		_, err := s.categoryRepo.FindByID(req.ReferenceCategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("reference category not found")
			}
			return nil, err
		}
	}

	reference.ReferenceCategoryID = req.ReferenceCategoryID
	reference.Name = req.Name
	if req.IsActive != nil {
		reference.IsActive = *req.IsActive
	}

	if err := s.referenceRepo.Update(reference); err != nil {
		return nil, err
	}

	// Reload with relations
	reference, err = s.referenceRepo.FindByID(reference.ID)
	if err != nil {
		return nil, err
	}

	response := model.ToReferenceResponse(reference)
	return &response, nil
}

func (s *referenceServiceImpl) DeleteReference(id int64) error {
	reference, err := s.referenceRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("reference not found")
		}
		return err
	}

	return s.referenceRepo.Delete(reference.ID)
}
