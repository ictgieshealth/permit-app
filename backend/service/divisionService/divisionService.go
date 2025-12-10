package divisionService

import (
	"errors"
	"permit-app/model"
	"permit-app/repo/divisionRepository"
)

type DivisionService interface {
	CreateDivision(req *model.DivisionRequest) (*model.DivisionResponse, error)
	GetDivisionByID(id int64) (*model.DivisionResponse, error)
	GetAllDivisions(filter *model.DivisionListRequest) ([]model.DivisionResponse, int64, error)
	UpdateDivision(id int64, req *model.DivisionUpdateRequest) (*model.DivisionResponse, error)
	DeleteDivision(id int64) error
}

type divisionService struct {
	repo divisionRepository.DivisionRepository
}

func NewDivisionService(repo divisionRepository.DivisionRepository) DivisionService {
	return &divisionService{repo: repo}
}

func (s *divisionService) CreateDivision(req *model.DivisionRequest) (*model.DivisionResponse, error) {
	existing, err := s.repo.FindByCodeAndDomainID(req.Code, req.DomainID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("division code already exists in this domain")
	}

	division := &model.Division{
		DomainID: req.DomainID,
		Code:     req.Code,
		Name:     req.Name,
	}

	err = s.repo.Create(division)
	if err != nil {
		return nil, err
	}

	// Reload with domain
	division, err = s.repo.FindByID(division.ID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(division), nil
}

func (s *divisionService) GetDivisionByID(id int64) (*model.DivisionResponse, error) {
	division, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(division), nil
}

func (s *divisionService) GetAllDivisions(filter *model.DivisionListRequest) ([]model.DivisionResponse, int64, error) {
	divisions, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.DivisionResponse
	for _, division := range divisions {
		responses = append(responses, *s.toResponse(&division))
	}

	return responses, total, nil
}

func (s *divisionService) UpdateDivision(id int64, req *model.DivisionUpdateRequest) (*model.DivisionResponse, error) {
	division, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.DomainID > 0 {
		division.DomainID = req.DomainID
	}

	if req.Code != "" && (req.Code != division.Code || req.DomainID > 0) {
		domainIDToCheck := division.DomainID
		if req.DomainID > 0 {
			domainIDToCheck = req.DomainID
		}
		existing, err := s.repo.FindByCodeAndDomainID(req.Code, domainIDToCheck)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("division code already exists in this domain")
		}
		division.Code = req.Code
	}

	if req.Name != "" {
		division.Name = req.Name
	}

	err = s.repo.Update(id, division)
	if err != nil {
		return nil, err
	}

	updatedDivision, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updatedDivision), nil
}

func (s *divisionService) DeleteDivision(id int64) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *divisionService) toResponse(division *model.Division) *model.DivisionResponse {
	resp := &model.DivisionResponse{
		ID:        division.ID,
		DomainID:  division.DomainID,
		Code:      division.Code,
		Name:      division.Name,
		CreatedAt: division.CreatedAt,
		UpdatedAt: division.UpdatedAt,
	}

	if division.Domain != nil {
		resp.Domain = &model.DomainResponse{
			ID:          division.Domain.ID,
			Code:        division.Domain.Code,
			Name:        division.Domain.Name,
			Description: division.Domain.Description,
			IsActive:    division.Domain.IsActive,
			CreatedAt:   division.Domain.CreatedAt,
			UpdatedAt:   division.Domain.UpdatedAt,
		}
	}

	return resp
}