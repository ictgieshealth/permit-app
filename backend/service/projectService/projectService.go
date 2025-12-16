package projectService

import (
	"errors"
	"permit-app/model"
	"permit-app/repo/projectRepository"
	"permit-app/repo/referenceRepository"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ProjectService interface {
	CreateProject(req model.ProjectRequest) (*model.ProjectResponse, error)
	GetProjects(req model.ProjectListRequest) ([]model.ProjectResponse, int64, error)
	GetProjectByID(id int64) (*model.ProjectResponse, error)
	GetProjectsByDomainID(domainID int64) ([]model.ProjectResponse, error)
	GetProjectsByUserID(userID int64) ([]model.ProjectResponse, error)
	UpdateProject(id int64, req model.ProjectUpdateRequest) (*model.ProjectResponse, error)
	DeleteProject(id int64) error
	ChangeProjectStatus(id int64, req model.ProjectStatusChangeRequest) (*model.ProjectResponse, error)
}

type projectServiceImpl struct {
	projectRepo   projectRepository.ProjectRepository
	referenceRepo referenceRepository.ReferenceRepository
}

func NewProjectService(projectRepo projectRepository.ProjectRepository, referenceRepo referenceRepository.ReferenceRepository) ProjectService {
	return &projectServiceImpl{
		projectRepo:   projectRepo,
		referenceRepo: referenceRepo,
	}
}

func (s *projectServiceImpl) CreateProject(req model.ProjectRequest) (*model.ProjectResponse, error) {
	// Generate code from name if not provided
	code := req.Code
	if code == "" {
		code = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
	}

	// Set default project status to "Pending" (26) if not provided
	projectStatusID := req.ProjectStatusID
	if projectStatusID == nil {
		defaultStatusID := int64(26) // Pending
		projectStatusID = &defaultStatusID
	}

	// Set default status to true if not provided
	status := req.Status
	if status == nil {
		defaultStatus := true
		status = &defaultStatus
	}

	project := &model.Project{
		DomainID:        req.DomainID,
		Name:            req.Name,
		Code:            code,
		Description:     req.Description,
		Status:          *status,
		ProjectStatusID: *projectStatusID,
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	// Assign users to project if provided
	if len(req.UserIDs) > 0 {
		if err := s.projectRepo.AssignUsers(project.ID, req.DomainID, req.UserIDs); err != nil {
			return nil, err
		}
	}

	// Fetch the created project with relations
	createdProject, err := s.projectRepo.FindByID(project.ID)
	if err != nil {
		return nil, err
	}

	response := model.ToProjectResponse(createdProject)
	return &response, nil
}

func (s *projectServiceImpl) GetProjects(req model.ProjectListRequest) ([]model.ProjectResponse, int64, error) {
	// Set default pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	filters := make(map[string]interface{})

	if req.DomainID > 0 {
		filters["domain_id"] = req.DomainID
	}
	if req.ProjectStatusID > 0 {
		filters["project_status_id"] = req.ProjectStatusID
	}
	if req.Name != "" {
		filters["name"] = req.Name
	}
	if req.Code != "" {
		filters["code"] = req.Code
	}
	if req.Status != nil {
		filters["status"] = *req.Status
	}

	projects, total, err := s.projectRepo.FindAll(filters, req.Page, req.Limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]model.ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = model.ToProjectResponse(&project)
	}

	return responses, total, nil
}

func (s *projectServiceImpl) GetProjectByID(id int64) (*model.ProjectResponse, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	response := model.ToProjectResponse(project)
	return &response, nil
}

func (s *projectServiceImpl) GetProjectsByDomainID(domainID int64) ([]model.ProjectResponse, error) {
	projects, err := s.projectRepo.FindByDomainID(domainID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = model.ToProjectResponse(&project)
	}

	return responses, nil
}

func (s *projectServiceImpl) GetProjectsByUserID(userID int64) ([]model.ProjectResponse, error) {
	projects, err := s.projectRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = model.ToProjectResponse(&project)
	}

	return responses, nil
}

func (s *projectServiceImpl) UpdateProject(id int64, req model.ProjectUpdateRequest) (*model.ProjectResponse, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	// Update fields
	project.Name = req.Name
	project.Description = req.Description

	if req.Status != nil {
		project.Status = *req.Status
	}

	if req.ProjectStatusID != nil {
		project.ProjectStatusID = *req.ProjectStatusID
	}

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	// Update user assignments if provided
	if req.UserIDs != nil {
		if err := s.projectRepo.AssignUsers(project.ID, project.DomainID, req.UserIDs); err != nil {
			return nil, err
		}
	}

	// Fetch updated project with relations
	updatedProject, err := s.projectRepo.FindByID(project.ID)
	if err != nil {
		return nil, err
	}

	response := model.ToProjectResponse(updatedProject)
	return &response, nil
}

func (s *projectServiceImpl) DeleteProject(id int64) error {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("project not found")
		}
		return err
	}

	// Check if project status is "Pending" (26) before deletion
	if project.ProjectStatusID != 26 {
		return errors.New("project cannot be deleted because it is not in pending status")
	}

	// Delete user assignments first
	if err := s.projectRepo.AssignUsers(project.ID, project.DomainID, []int64{}); err != nil {
		return err
	}

	return s.projectRepo.Delete(id)
}

func (s *projectServiceImpl) ChangeProjectStatus(id int64, req model.ProjectStatusChangeRequest) (*model.ProjectResponse, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	// Verify the status exists
	_, err = s.referenceRepo.FindByID(req.StatusID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project status not found")
		}
		return nil, err
	}

	// Prepare updates
	updates := map[string]interface{}{
		"project_status_id": req.StatusID,
	}

	// Set started_at when status changes to "On Progress" (28)
	if req.StatusID == 28 && project.StartedAt == nil {
		now := time.Now()
		updates["started_at"] = now
	}

	// Set finished_at when status changes to "Done" (29)
	if req.StatusID == 29 && project.FinishedAt == nil {
		now := time.Now()
		updates["finished_at"] = now
	}

	// Update only specific fields without affecting associations
	if err := s.projectRepo.UpdateFields(id, updates); err != nil {
		return nil, err
	}

	// Fetch updated project with relations
	updatedProject, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := model.ToProjectResponse(updatedProject)
	return &response, nil
}
