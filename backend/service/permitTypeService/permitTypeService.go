package permitTypeService

import (
	"permit-app/model"
	"permit-app/repo/permitTypeRepository"
)

type PermitTypeService interface {
	CreatePermitType(req *model.PermitTypeRequest) (*model.PermitTypeResponse, error)
	GetPermitTypeByID(id int64) (*model.PermitTypeResponse, error)
	GetAllPermitTypes(filter *model.PermitTypeListRequest) ([]model.PermitTypeResponse, int64, error)
	UpdatePermitType(id int64, req *model.PermitTypeUpdateRequest) (*model.PermitTypeResponse, error)
	DeletePermitType(id int64) error
}

type permitTypeService struct {
	repo permitTypeRepository.PermitTypeRepository
}

func NewPermitTypeService(repo permitTypeRepository.PermitTypeRepository) PermitTypeService {
	return &permitTypeService{repo: repo}
}

func (s *permitTypeService) CreatePermitType(req *model.PermitTypeRequest) (*model.PermitTypeResponse, error) {
	permitType := &model.PermitType{
		DivisionID:             req.DivisionID,
		Name:                   req.Name,
		RiskPoint:              req.RiskPoint,
		DefaultApplicationType: req.DefaultApplicationType,
		DefaultValidityPeriod:  req.DefaultValidityPeriod,
		Notes:                  req.Notes,
	}

	err := s.repo.Create(permitType)
	if err != nil {
		return nil, err
	}

	created, err := s.repo.FindByID(permitType.ID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(created), nil
}

func (s *permitTypeService) GetPermitTypeByID(id int64) (*model.PermitTypeResponse, error) {
	permitType, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(permitType), nil
}

func (s *permitTypeService) GetAllPermitTypes(filter *model.PermitTypeListRequest) ([]model.PermitTypeResponse, int64, error) {
	permitTypes, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.PermitTypeResponse
	for _, permitType := range permitTypes {
		responses = append(responses, *s.toResponse(&permitType))
	}

	return responses, total, nil
}

func (s *permitTypeService) UpdatePermitType(id int64, req *model.PermitTypeUpdateRequest) (*model.PermitTypeResponse, error) {
	permitType, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.DivisionID != nil {
		permitType.DivisionID = req.DivisionID
	}
	if req.Name != "" {
		permitType.Name = req.Name
	}
	if req.RiskPoint != nil {
		permitType.RiskPoint = req.RiskPoint
	}
	if req.DefaultApplicationType != nil {
		permitType.DefaultApplicationType = req.DefaultApplicationType
	}
	if req.DefaultValidityPeriod != nil {
		permitType.DefaultValidityPeriod = req.DefaultValidityPeriod
	}
	if req.Notes != nil {
		permitType.Notes = req.Notes
	}

	err = s.repo.Update(id, permitType)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updated), nil
}

func (s *permitTypeService) DeletePermitType(id int64) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *permitTypeService) toResponse(permitType *model.PermitType) *model.PermitTypeResponse {
	response := &model.PermitTypeResponse{
		ID:                     permitType.ID,
		DivisionID:             permitType.DivisionID,
		Name:                   permitType.Name,
		RiskPoint:              permitType.RiskPoint,
		DefaultApplicationType: permitType.DefaultApplicationType,
		DefaultValidityPeriod:  permitType.DefaultValidityPeriod,
		Notes:                  permitType.Notes,
		CreatedAt:              permitType.CreatedAt,
		UpdatedAt:              permitType.UpdatedAt,
	}

	if permitType.Division != nil && permitType.Division.ID > 0 {
		response.Division = &model.DivisionResponse{
			ID:        permitType.Division.ID,
			DomainID:  permitType.Division.DomainID,
			Code:      permitType.Division.Code,
			Name:      permitType.Division.Name,
			CreatedAt: permitType.Division.CreatedAt,
			UpdatedAt: permitType.Division.UpdatedAt,
		}

		if permitType.Division.Domain != nil {
			response.Division.Domain = &model.DomainResponse{
				ID:          permitType.Division.Domain.ID,
				Code:        permitType.Division.Domain.Code,
				Name:        permitType.Division.Domain.Name,
				Description: permitType.Division.Domain.Description,
				IsActive:    permitType.Division.Domain.IsActive,
				CreatedAt:   permitType.Division.Domain.CreatedAt,
				UpdatedAt:   permitType.Division.Domain.UpdatedAt,
			}
		}
	}

	return response
}