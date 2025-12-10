package menuService

import (
	"errors"
	"permit-app/model"
	"permit-app/repo/menuRepository"

	"gorm.io/gorm"
)

type MenuService interface {
	CreateMenu(req *model.MenuRequest) (*model.Menu, error)
	UpdateMenu(id uint, req *model.MenuRequest) (*model.Menu, error)
	DeleteMenu(id uint) error
	GetMenuByID(id uint) (*model.Menu, error)
	GetAllMenus(filter *model.MenuFilter) ([]model.Menu, int64, error)
	GetUserMenus(roleID uint) ([]model.Menu, error)
	AssignRolesToMenu(menuID uint, roleIDs []uint) error
}

type menuService struct {
	menuRepo menuRepository.MenuRepository
}

func NewMenuService(menuRepo menuRepository.MenuRepository) MenuService {
	return &menuService{menuRepo: menuRepo}
}

func (s *menuService) CreateMenu(req *model.MenuRequest) (*model.Menu, error) {
	// Check if path already exists
	existing, err := s.menuRepo.FindByPath(req.Path)
	if err == nil && existing != nil {
		return nil, errors.New("menu path already exists")
	}

	// Validate parent_id if provided
	if req.ParentID != nil {
		parent, err := s.menuRepo.FindByID(*req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("parent menu not found")
			}
			return nil, err
		}
		// Check for circular reference
		if parent.ParentID != nil && *parent.ParentID == *req.ParentID {
			return nil, errors.New("circular reference detected")
		}
	}

	// Validate role IDs
	if len(req.RoleIDs) == 0 {
		return nil, errors.New("at least one role must be assigned")
	}

	menu := &model.Menu{
		Name:       req.Name,
		Path:       req.Path,
		Icon:       req.Icon,
		ParentID:   req.ParentID,
		OrderIndex: req.OrderIndex,
		IsActive:   true,
	}

	// Create menu
	if err := s.menuRepo.Create(menu); err != nil {
		return nil, err
	}

	// Assign roles
	if err := s.menuRepo.AssignRoles(menu.ID, req.RoleIDs); err != nil {
		// Rollback: delete the created menu
		s.menuRepo.Delete(menu.ID)
		return nil, err
	}

	// Fetch the created menu with roles
	return s.menuRepo.FindByID(menu.ID)
}

func (s *menuService) UpdateMenu(id uint, req *model.MenuRequest) (*model.Menu, error) {
	// Check if menu exists
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("menu not found")
		}
		return nil, err
	}

	// Check if path is being changed and if new path already exists
	if req.Path != menu.Path {
		existing, err := s.menuRepo.FindByPath(req.Path)
		if err == nil && existing != nil && existing.ID != id {
			return nil, errors.New("menu path already exists")
		}
	}

	// Validate parent_id if provided
	if req.ParentID != nil {
		// Can't set itself as parent
		if *req.ParentID == id {
			return nil, errors.New("menu cannot be its own parent")
		}
		parent, err := s.menuRepo.FindByID(*req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("parent menu not found")
			}
			return nil, err
		}
		// Check for circular reference
		if parent.ParentID != nil && *parent.ParentID == id {
			return nil, errors.New("circular reference detected")
		}
	}

	// Validate role IDs
	if len(req.RoleIDs) == 0 {
		return nil, errors.New("at least one role must be assigned")
	}

	// Update menu fields
	menu.Name = req.Name
	menu.Path = req.Path
	menu.Icon = req.Icon
	menu.ParentID = req.ParentID
	menu.OrderIndex = req.OrderIndex

	// Update menu
	if err := s.menuRepo.Update(menu); err != nil {
		return nil, err
	}

	// Update role assignments
	if err := s.menuRepo.AssignRoles(menu.ID, req.RoleIDs); err != nil {
		return nil, err
	}

	// Fetch the updated menu with roles
	return s.menuRepo.FindByID(menu.ID)
}

func (s *menuService) DeleteMenu(id uint) error {
	// Check if menu exists
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("menu not found")
		}
		return err
	}

	// Remove role assignments first
	if err := s.menuRepo.RemoveRoles(menu.ID); err != nil {
		return err
	}

	// Delete menu (will cascade delete children and menu_roles)
	return s.menuRepo.Delete(id)
}

func (s *menuService) GetMenuByID(id uint) (*model.Menu, error) {
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("menu not found")
		}
		return nil, err
	}
	return menu, nil
}

func (s *menuService) GetAllMenus(filter *model.MenuFilter) ([]model.Menu, int64, error) {
	// Set default values
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	if filter.Page == 0 {
		filter.Page = 1
	}

	return s.menuRepo.FindAll(filter)
}

func (s *menuService) GetUserMenus(roleID uint) ([]model.Menu, error) {
	return s.menuRepo.FindByUserRole(roleID)
}

func (s *menuService) AssignRolesToMenu(menuID uint, roleIDs []uint) error {
	// Check if menu exists
	_, err := s.menuRepo.FindByID(menuID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("menu not found")
		}
		return err
	}

	// Validate role IDs
	if len(roleIDs) == 0 {
		return errors.New("at least one role must be assigned")
	}

	return s.menuRepo.AssignRoles(menuID, roleIDs)
}
