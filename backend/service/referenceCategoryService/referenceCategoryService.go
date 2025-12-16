package referenceCategoryService

import (
	"errors"
	"permit-app/model"
	"permit-app/repo/moduleRepository"
	"permit-app/repo/referenceCategoryRepository"

	"gorm.io/gorm"
)

type ReferenceCategoryService interface {
	CreateReferenceCategory(req model.ReferenceCategoryRequest) (*model.ReferenceCategoryResponse, error)
	GetReferenceCategories(req model.ReferenceCategoryListRequest) ([]model.ReferenceCategoryResponse, int64, error)
	GetReferenceCategoryByID(id int64) (*model.ReferenceCategoryResponse, error)
	GetCategoriesByModuleID(moduleID int64) ([]model.ReferenceCategoryResponse, error)
	UpdateReferenceCategory(id int64, req model.ReferenceCategoryRequest) (*model.ReferenceCategoryResponse, error)
	DeleteReferenceCategory(id int64) error
}

type referenceCategoryServiceImpl struct {
	categoryRepo referenceCategoryRepository.ReferenceCategoryRepository
	moduleRepo   moduleRepository.ModuleRepository
}

func NewReferenceCategoryService(
	categoryRepo referenceCategoryRepository.ReferenceCategoryRepository,
	moduleRepo moduleRepository.ModuleRepository,
) ReferenceCategoryService {
	return &referenceCategoryServiceImpl{
		categoryRepo: categoryRepo,
		moduleRepo:   moduleRepo,
	}
}

func (s *referenceCategoryServiceImpl) CreateReferenceCategory(req model.ReferenceCategoryRequest) (*model.ReferenceCategoryResponse, error) {
	// Validate module exists
	_, err := s.moduleRepo.FindByID(req.ModuleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, err
	}

	category := &model.ReferenceCategory{
		ModuleID: req.ModuleID,
		Name:     req.Name,
		IsActive: true,
	}

	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	// Reload with relations
	category, err = s.categoryRepo.FindByID(category.ID)
	if err != nil {
		return nil, err
	}

	response := model.ToReferenceCategoryResponse(category)
	return &response, nil
}

func (s *referenceCategoryServiceImpl) GetReferenceCategories(req model.ReferenceCategoryListRequest) ([]model.ReferenceCategoryResponse, int64, error) {
	filters := make(map[string]interface{})

	if req.ModuleID > 0 {
		filters["module_id"] = req.ModuleID
	}
	if req.Name != "" {
		filters["name"] = req.Name
	}
	if req.IsActive != nil {
		filters["is_active"] = *req.IsActive
	}

	categories, total, err := s.categoryRepo.FindAll(filters, req.Page, req.Limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response
	responses := make([]model.ReferenceCategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = model.ToReferenceCategoryResponse(&category)
	}

	return responses, total, nil
}

func (s *referenceCategoryServiceImpl) GetReferenceCategoryByID(id int64) (*model.ReferenceCategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reference category not found")
		}
		return nil, err
	}

	response := model.ToReferenceCategoryResponse(category)
	return &response, nil
}

func (s *referenceCategoryServiceImpl) GetCategoriesByModuleID(moduleID int64) ([]model.ReferenceCategoryResponse, error) {
	categories, err := s.categoryRepo.FindByModuleID(moduleID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ReferenceCategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = model.ToReferenceCategoryResponse(&category)
	}

	return responses, nil
}

func (s *referenceCategoryServiceImpl) UpdateReferenceCategory(id int64, req model.ReferenceCategoryRequest) (*model.ReferenceCategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reference category not found")
		}
		return nil, err
	}

	// Validate module exists if changed
	if req.ModuleID != category.ModuleID {
		_, err := s.moduleRepo.FindByID(req.ModuleID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("module not found")
			}
			return nil, err
		}
	}

	category.ModuleID = req.ModuleID
	category.Name = req.Name
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	// Reload with relations
	category, err = s.categoryRepo.FindByID(category.ID)
	if err != nil {
		return nil, err
	}

	response := model.ToReferenceCategoryResponse(category)
	return &response, nil
}

func (s *referenceCategoryServiceImpl) DeleteReferenceCategory(id int64) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("reference category not found")
		}
		return err
	}

	return s.categoryRepo.Delete(category.ID)
}
