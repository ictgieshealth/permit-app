package moduleService

import (
	"errors"
	"permit-app/model"
	"permit-app/repo/moduleRepository"

	"gorm.io/gorm"
)

type ModuleService interface {
	CreateModule(req model.ModuleRequest) (*model.ModuleResponse, error)
	GetModules(req model.ModuleListRequest) ([]model.ModuleResponse, int64, error)
	GetModuleByID(id int64) (*model.ModuleResponse, error)
	UpdateModule(id int64, req model.ModuleRequest) (*model.ModuleResponse, error)
	DeleteModule(id int64) error
}

type moduleServiceImpl struct {
	moduleRepo moduleRepository.ModuleRepository
}

func NewModuleService(moduleRepo moduleRepository.ModuleRepository) ModuleService {
	return &moduleServiceImpl{
		moduleRepo: moduleRepo,
	}
}

func (s *moduleServiceImpl) CreateModule(req model.ModuleRequest) (*model.ModuleResponse, error) {
	// Check if code already exists
	existing, _ := s.moduleRepo.FindByCode(req.Code)
	if existing != nil {
		return nil, errors.New("module code already exists")
	}

	module := &model.Module{
		Code:     req.Code,
		Name:     req.Name,
		IsActive: true,
	}

	if req.IsActive != nil {
		module.IsActive = *req.IsActive
	}

	if err := s.moduleRepo.Create(module); err != nil {
		return nil, err
	}

	response := model.ToModuleResponse(module)
	return &response, nil
}

func (s *moduleServiceImpl) GetModules(req model.ModuleListRequest) ([]model.ModuleResponse, int64, error) {
	filters := make(map[string]interface{})

	if req.Code != "" {
		filters["code"] = req.Code
	}
	if req.Name != "" {
		filters["name"] = req.Name
	}
	if req.IsActive != nil {
		filters["is_active"] = *req.IsActive
	}

	modules, total, err := s.moduleRepo.FindAll(filters, req.Page, req.Limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response
	responses := make([]model.ModuleResponse, len(modules))
	for i, module := range modules {
		responses[i] = model.ToModuleResponse(&module)
	}

	return responses, total, nil
}

func (s *moduleServiceImpl) GetModuleByID(id int64) (*model.ModuleResponse, error) {
	module, err := s.moduleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, err
	}

	response := model.ToModuleResponse(module)
	return &response, nil
}

func (s *moduleServiceImpl) UpdateModule(id int64, req model.ModuleRequest) (*model.ModuleResponse, error) {
	module, err := s.moduleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, err
	}

	// Check if new code conflicts with existing
	if req.Code != module.Code {
		existing, _ := s.moduleRepo.FindByCode(req.Code)
		if existing != nil {
			return nil, errors.New("module code already exists")
		}
	}

	module.Code = req.Code
	module.Name = req.Name
	if req.IsActive != nil {
		module.IsActive = *req.IsActive
	}

	if err := s.moduleRepo.Update(module); err != nil {
		return nil, err
	}

	response := model.ToModuleResponse(module)
	return &response, nil
}

func (s *moduleServiceImpl) DeleteModule(id int64) error {
	module, err := s.moduleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("module not found")
		}
		return err
	}

	return s.moduleRepo.Delete(module.ID)
}
