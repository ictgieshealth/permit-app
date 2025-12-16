package projectRepository

import (
	"permit-app/model"
	"strings"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *model.Project) error
	FindAll(filters map[string]interface{}, page, limit int) ([]model.Project, int64, error)
	FindByID(id int64) (*model.Project, error)
	FindByDomainID(domainID int64) ([]model.Project, error)
	FindByUserID(userID int64) ([]model.Project, error)
	Update(project *model.Project) error
	UpdateFields(id int64, fields map[string]interface{}) error
	Delete(id int64) error
	AssignUsers(projectID, domainID int64, userIDs []int64) error
	RemoveUsers(projectID int64, userIDs []int64) error
	GetProjectUsers(projectID int64) ([]model.User, error)
}

type projectRepositoryImpl struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepositoryImpl{db: db}
}

func (r *projectRepositoryImpl) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepositoryImpl) FindAll(filters map[string]interface{}, page, limit int) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	query := r.db.Model(&model.Project{}).
		Preload("Domain").
		Preload("ProjectStatus").
		Preload("ProjectStatus.ReferenceCategory").
		Preload("ProjectStatus.ReferenceCategory.Module").
		Preload("Users")

	// Apply filters
	if domainID, ok := filters["domain_id"].(int64); ok && domainID > 0 {
		query = query.Where("domain_id = ?", domainID)
	}
	if projectStatusID, ok := filters["project_status_id"].(int64); ok && projectStatusID > 0 {
		query = query.Where("project_status_id = ?", projectStatusID)
	}
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%")
	}
	if code, ok := filters["code"].(string); ok && code != "" {
		query = query.Where("LOWER(code) LIKE ?", "%"+strings.ToLower(code)+"%")
	}
	if status, ok := filters["status"].(bool); ok {
		query = query.Where("status = ?", status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *projectRepositoryImpl) FindByID(id int64) (*model.Project, error) {
	var project model.Project
	err := r.db.
		Preload("Domain").
		Preload("ProjectStatus").
		Preload("ProjectStatus.ReferenceCategory").
		Preload("ProjectStatus.ReferenceCategory.Module").
		Preload("Users").
		First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepositoryImpl) FindByDomainID(domainID int64) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.
		Preload("Domain").
		Preload("ProjectStatus").
		Preload("Users").
		Where("domain_id = ?", domainID).
		Order("created_at DESC").
		Find(&projects).Error
	return projects, err
}

func (r *projectRepositoryImpl) FindByUserID(userID int64) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.
		Preload("Domain").
		Preload("ProjectStatus").
		Preload("Users").
		Joins("JOIN user_domain_projects ON user_domain_projects.project_id = projects.id").
		Where("user_domain_projects.user_id = ?", userID).
		Order("projects.created_at DESC").
		Find(&projects).Error
	return projects, err
}

func (r *projectRepositoryImpl) Update(project *model.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepositoryImpl) UpdateFields(id int64, fields map[string]interface{}) error {
	return r.db.Model(&model.Project{}).Where("id = ?", id).Updates(fields).Error
}

func (r *projectRepositoryImpl) Delete(id int64) error {
	return r.db.Delete(&model.Project{}, id).Error
}

func (r *projectRepositoryImpl) AssignUsers(projectID, domainID int64, userIDs []int64) error {
	// First, remove existing assignments
	if err := r.db.Where("project_id = ? AND domain_id = ?", projectID, domainID).Delete(&model.UserDomainProject{}).Error; err != nil {
		return err
	}

	// Then, create new assignments
	if len(userIDs) > 0 {
		assignments := make([]model.UserDomainProject, len(userIDs))
		for i, userID := range userIDs {
			assignments[i] = model.UserDomainProject{
				UserID:    userID,
				DomainID:  domainID,
				ProjectID: projectID,
			}
		}
		return r.db.Create(&assignments).Error
	}

	return nil
}

func (r *projectRepositoryImpl) RemoveUsers(projectID int64, userIDs []int64) error {
	return r.db.Where("project_id = ? AND user_id IN ?", projectID, userIDs).Delete(&model.UserDomainProject{}).Error
}

func (r *projectRepositoryImpl) GetProjectUsers(projectID int64) ([]model.User, error) {
	var users []model.User
	err := r.db.
		Joins("JOIN user_domain_projects ON user_domain_projects.user_id = users.id").
		Where("user_domain_projects.project_id = ?", projectID).
		Find(&users).Error
	return users, err
}
