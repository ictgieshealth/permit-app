package permitService

import (
	"errors"
	"mime/multipart"
	"permit-app/helper"
	"permit-app/model"
	"permit-app/repo/permitRepository"
)

type PermitService interface {
	CreatePermit(req *model.PermitRequest) (*model.PermitResponse, error)
	GetPermitByID(id int64) (*model.PermitResponse, error)
	GetAllPermits(filter *model.PermitListRequest) ([]model.PermitResponse, int64, error)
	UpdatePermit(id int64, req *model.PermitUpdateRequest) (*model.PermitResponse, error)
	DeletePermit(id int64) error
	HandleFileUpload(id int64, file *multipart.FileHeader) (*model.PermitResponse, error)
	SearchPermits(query string, filter *model.PermitListRequest) ([]model.PermitResponse, int64, error)
}

type permitService struct {
	repo permitRepository.PermitRepository
}

func NewPermitService(repo permitRepository.PermitRepository) PermitService {
	return &permitService{repo: repo}
}

func (s *permitService) CreatePermit(req *model.PermitRequest) (*model.PermitResponse, error) {
	existing, err := s.repo.FindByPermitNoAndDomainID(req.PermitNo, req.DomainID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("permit number already exists in this domain")
	}

	permit := &model.Permit{
		DomainID:               req.DomainID,
		DivisionID:             req.DivisionID,
		PermitTypeID:           req.PermitTypeID,
		Name:                   req.Name,
		ApplicationType:        req.ApplicationType,
		PermitNo:               req.PermitNo,
		EffectiveDate:          req.EffectiveDate.Time,
		ExpiryDate:             req.ExpiryDate.Time,
		EffectiveTerm:          req.EffectiveTerm,
		ResponsiblePersonID:    req.ResponsiblePersonID,
		ResponsibleDocPersonID: req.ResponsibleDocPersonID,
		DocName:                req.DocName,
		DocNumber:              req.DocNumber,
		Status:                 req.Status,
	}

	err = s.repo.Create(permit)
	if err != nil {
		return nil, err
	}

	created, err := s.repo.FindByID(permit.ID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(created), nil
}

func (s *permitService) GetPermitByID(id int64) (*model.PermitResponse, error) {
	permit, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(permit), nil
}

func (s *permitService) GetAllPermits(filter *model.PermitListRequest) ([]model.PermitResponse, int64, error) {
	permits, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.PermitResponse
	for _, permit := range permits {
		responses = append(responses, *s.toResponse(&permit))
	}

	return responses, total, nil
}

func (s *permitService) UpdatePermit(id int64, req *model.PermitUpdateRequest) (*model.PermitResponse, error) {
	permit, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.DomainID > 0 {
		permit.DomainID = req.DomainID
	}

	if req.DivisionID != nil {
		permit.DivisionID = req.DivisionID
	}

	if req.PermitNo != "" && (req.PermitNo != permit.PermitNo || req.DomainID > 0) {
		domainIDToCheck := permit.DomainID
		if req.DomainID > 0 {
			domainIDToCheck = req.DomainID
		}
		existing, err := s.repo.FindByPermitNoAndDomainID(req.PermitNo, domainIDToCheck)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("permit number already exists in this domain")
		}
		permit.PermitNo = req.PermitNo
	}

	if req.PermitTypeID > 0 {
		permit.PermitTypeID = req.PermitTypeID
	}
	if req.Name != "" {
		permit.Name = req.Name
	}
	if req.ApplicationType != "" {
		permit.ApplicationType = req.ApplicationType
	}
	if !req.EffectiveDate.IsZero() {
		permit.EffectiveDate = req.EffectiveDate.Time
	}
	if !req.ExpiryDate.IsZero() {
		permit.ExpiryDate = req.ExpiryDate.Time
	}
	if req.EffectiveTerm != nil {
		permit.EffectiveTerm = req.EffectiveTerm
	}
	if req.ResponsiblePersonID != nil {
		permit.ResponsiblePersonID = req.ResponsiblePersonID
	}
	if req.ResponsibleDocPersonID != nil {
		permit.ResponsibleDocPersonID = req.ResponsibleDocPersonID
	}
	if req.DocName != nil {
		permit.DocName = req.DocName
	}
	if req.DocNumber != nil {
		permit.DocNumber = req.DocNumber
	}
	if req.Status != "" {
		permit.Status = req.Status
	}

	err = s.repo.Update(id, permit)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updated), nil
}

func (s *permitService) DeletePermit(id int64) error {
	permit, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Delete associated file if exists
	if permit.DocFilePath != nil && *permit.DocFilePath != "" {
		helper.DeleteFile(*permit.DocFilePath)
	}

	return s.repo.Delete(id)
}

func (s *permitService) HandleFileUpload(id int64, file *multipart.FileHeader) (*model.PermitResponse, error) {
	// Get existing permit
	permit, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Delete old file if exists
	if permit.DocFilePath != nil && *permit.DocFilePath != "" {
		helper.DeleteFile(*permit.DocFilePath)
	}

	// Save new file
	filePath, err := helper.SaveFile(file, "permits")
	if err != nil {
		return nil, err
	}

	// Update permit with file information
	fileName := file.Filename
	fileSize := file.Size
	fileType := helper.GetMimeType(file.Filename)

	permit.DocFileName = &fileName
	permit.DocFilePath = &filePath
	permit.DocFileSize = &fileSize
	permit.DocFileType = &fileType

	err = s.repo.Update(id, permit)
	if err != nil {
		// If database update fails, delete the uploaded file
		helper.DeleteFile(filePath)
		return nil, err
	}

	// Fetch updated permit with all relations
	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updated), nil
}

func (s *permitService) toResponse(permit *model.Permit) *model.PermitResponse {
	resp := &model.PermitResponse{
		ID:                     permit.ID,
		DomainID:               permit.DomainID,
		DivisionID:             permit.DivisionID,
		PermitTypeID:           permit.PermitTypeID,
		Name:                   permit.Name,
		ApplicationType:        permit.ApplicationType,
		PermitNo:               permit.PermitNo,
		EffectiveDate:          permit.EffectiveDate,
		ExpiryDate:             permit.ExpiryDate,
		EffectiveTerm:          permit.EffectiveTerm,
		ResponsiblePersonID:    permit.ResponsiblePersonID,
		ResponsibleDocPersonID: permit.ResponsibleDocPersonID,
		DocName:                permit.DocName,
		DocNumber:              permit.DocNumber,
		DocFileName:            permit.DocFileName,
		DocFilePath:            permit.DocFilePath,
		DocFileSize:            permit.DocFileSize,
		DocFileType:            permit.DocFileType,
		Status:                 permit.Status,
		CreatedAt:              permit.CreatedAt,
		UpdatedAt:              permit.UpdatedAt,
	}

	if permit.Domain != nil {
		resp.Domain = &model.DomainResponse{
			ID:          permit.Domain.ID,
			Code:        permit.Domain.Code,
			Name:        permit.Domain.Name,
			Description: permit.Domain.Description,
			IsActive:    permit.Domain.IsActive,
			CreatedAt:   permit.Domain.CreatedAt,
			UpdatedAt:   permit.Domain.UpdatedAt,
		}
	}

	if permit.Division != nil {
		resp.Division = &model.DivisionResponse{
			ID:        permit.Division.ID,
			DomainID:  permit.Division.DomainID,
			Code:      permit.Division.Code,
			Name:      permit.Division.Name,
			CreatedAt: permit.Division.CreatedAt,
			UpdatedAt: permit.Division.UpdatedAt,
		}

		if permit.Division.Domain != nil {
			resp.Division.Domain = &model.DomainResponse{
				ID:          permit.Division.Domain.ID,
				Code:        permit.Division.Domain.Code,
				Name:        permit.Division.Domain.Name,
				Description: permit.Division.Domain.Description,
				IsActive:    permit.Division.Domain.IsActive,
				CreatedAt:   permit.Division.Domain.CreatedAt,
				UpdatedAt:   permit.Division.Domain.UpdatedAt,
			}
		}
	}

	if permit.PermitType != nil {
		resp.PermitType = &model.PermitTypeResponse{
			ID:                     permit.PermitType.ID,
			DivisionID:             permit.PermitType.DivisionID,
			Name:                   permit.PermitType.Name,
			RiskPoint:              permit.PermitType.RiskPoint,
			DefaultApplicationType: permit.PermitType.DefaultApplicationType,
			DefaultValidityPeriod:  permit.PermitType.DefaultValidityPeriod,
			Notes:                  permit.PermitType.Notes,
			CreatedAt:              permit.PermitType.CreatedAt,
			UpdatedAt:              permit.PermitType.UpdatedAt,
		}

		if permit.PermitType.Division != nil {
			resp.PermitType.Division = &model.DivisionResponse{
				ID:        permit.PermitType.Division.ID,
				DomainID:  permit.PermitType.Division.DomainID,
				Code:      permit.PermitType.Division.Code,
				Name:      permit.PermitType.Division.Name,
				CreatedAt: permit.PermitType.Division.CreatedAt,
				UpdatedAt: permit.PermitType.Division.UpdatedAt,
			}
		}
	}

	if permit.ResponsiblePerson != nil {
		resp.ResponsiblePerson = &model.UserResponse{
			ID:       permit.ResponsiblePerson.ID,
			Username: permit.ResponsiblePerson.Username,
			Email:    permit.ResponsiblePerson.Email,
			FullName: permit.ResponsiblePerson.FullName,
			IsActive: permit.ResponsiblePerson.IsActive,
		}
	}

	if permit.ResponsibleDocPerson != nil {
		resp.ResponsibleDocPerson = &model.UserResponse{
			ID:       permit.ResponsibleDocPerson.ID,
			Username: permit.ResponsibleDocPerson.Username,
			Email:    permit.ResponsibleDocPerson.Email,
			FullName: permit.ResponsibleDocPerson.FullName,
			IsActive: permit.ResponsibleDocPerson.IsActive,
		}
	}

	return resp
}

func (s *permitService) SearchPermits(query string, filter *model.PermitListRequest) ([]model.PermitResponse, int64, error) {
	permits, total, err := s.repo.Search(query, filter)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.PermitResponse
	for _, permit := range permits {
		responses = append(responses, *s.toResponse(&permit))
	}

	return responses, total, nil
}
