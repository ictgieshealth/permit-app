package domainService

import (
	"errors"
	"permit-app/model"
	"permit-app/repo/domainRepository"
)

type DomainService interface {
	CreateDomain(req *model.DomainRequest) (*model.DomainResponse, error)
	GetDomainByID(id int64) (*model.DomainResponse, error)
	GetAllDomains(filter *model.DomainListRequest) ([]model.DomainResponse, int64, error)
	UpdateDomain(id int64, req *model.DomainUpdateRequest) (*model.DomainResponse, error)
	DeleteDomain(id int64) error
}

type domainService struct {
	repo domainRepository.DomainRepository
}

func NewDomainService(repo domainRepository.DomainRepository) DomainService {
	return &domainService{repo: repo}
}

func (s *domainService) CreateDomain(req *model.DomainRequest) (*model.DomainResponse, error) {
	existing, err := s.repo.FindByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("domain code already exists")
	}

	domain := &model.Domain{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	if req.IsActive != nil {
		domain.IsActive = *req.IsActive
	}

	err = s.repo.Create(domain)
	if err != nil {
		return nil, err
	}

	return s.toResponse(domain), nil
}

func (s *domainService) GetDomainByID(id int64) (*model.DomainResponse, error) {
	domain, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(domain), nil
}

func (s *domainService) GetAllDomains(filter *model.DomainListRequest) ([]model.DomainResponse, int64, error) {
	domains, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.DomainResponse
	for _, domain := range domains {
		responses = append(responses, *s.toResponse(&domain))
	}

	return responses, total, nil
}

func (s *domainService) UpdateDomain(id int64, req *model.DomainUpdateRequest) (*model.DomainResponse, error) {
	domain, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Code != "" && req.Code != domain.Code {
		existing, err := s.repo.FindByCode(req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("domain code already exists")
		}
		domain.Code = req.Code
	}

	if req.Name != "" {
		domain.Name = req.Name
	}

	if req.Description != nil {
		domain.Description = req.Description
	}

	if req.IsActive != nil {
		domain.IsActive = *req.IsActive
	}

	err = s.repo.Update(id, domain)
	if err != nil {
		return nil, err
	}

	updatedDomain, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updatedDomain), nil
}

func (s *domainService) DeleteDomain(id int64) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *domainService) toResponse(domain *model.Domain) *model.DomainResponse {
	return &model.DomainResponse{
		ID:          domain.ID,
		Code:        domain.Code,
		Name:        domain.Name,
		Description: domain.Description,
		IsActive:    domain.IsActive,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}
