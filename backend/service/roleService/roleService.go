package roleService

import (
	"errors"
	"permit-app/model"
	"permit-app/repo/roleRepository"
)

type RoleService interface {
	CreateRole(req *model.RoleRequest) (*model.RoleResponse, error)
	GetRoleByID(id int64) (*model.RoleResponse, error)
	GetAllRoles(filter *model.RoleListRequest) ([]model.RoleResponse, int64, error)
	UpdateRole(id int64, req *model.RoleUpdateRequest) (*model.RoleResponse, error)
	DeleteRole(id int64) error
}

type roleService struct {
	repo roleRepository.RoleRepository
}

func NewRoleService(repo roleRepository.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) CreateRole(req *model.RoleRequest) (*model.RoleResponse, error) {
	// Check if code already exists
	existingCode, err := s.repo.FindByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if existingCode != nil {
		return nil, errors.New("role code already exists")
	}

	// Check if name already exists
	existing, err := s.repo.FindByName(req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("role name already exists")
	}

	role := &model.Role{
		Code:        req.Code,
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
	}

	err = s.repo.Create(role)
	if err != nil {
		return nil, err
	}

	return s.toResponse(role), nil
}

func (s *roleService) GetRoleByID(id int64) (*model.RoleResponse, error) {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(role), nil
}

func (s *roleService) GetAllRoles(filter *model.RoleListRequest) ([]model.RoleResponse, int64, error) {
	roles, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.RoleResponse
	for _, role := range roles {
		responses = append(responses, *s.toResponse(&role))
	}

	return responses, total, nil
}

func (s *roleService) UpdateRole(id int64, req *model.RoleUpdateRequest) (*model.RoleResponse, error) {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check code uniqueness if changed
	if req.Code != "" && req.Code != role.Code {
		existing, err := s.repo.FindByCode(req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("role code already exists")
		}
		role.Code = req.Code
	}

	// Check name uniqueness if changed
	if req.Name != "" && req.Name != role.Name {
		existing, err := s.repo.FindByName(req.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("role name already exists")
		}
		role.Name = req.Name
	}

	if req.Category != "" {
		role.Category = req.Category
	}

	if req.Description != nil {
		role.Description = req.Description
	}

	err = s.repo.Update(id, role)
	if err != nil {
		return nil, err
	}

	updatedRole, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updatedRole), nil
}

func (s *roleService) DeleteRole(id int64) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *roleService) toResponse(role *model.Role) *model.RoleResponse {
	return &model.RoleResponse{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Category:    role.Category,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
