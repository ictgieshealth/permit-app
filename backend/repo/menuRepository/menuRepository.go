package menuRepository

import (
	"permit-app/model"

	"gorm.io/gorm"
)

type MenuRepository interface {
	Create(menu *model.Menu) error
	Update(menu *model.Menu) error
	Delete(id uint) error
	FindByID(id uint) (*model.Menu, error)
	FindAll(filter *model.MenuFilter) ([]model.Menu, int64, error)
	FindByUserRole(roleID uint) ([]model.Menu, error)
	AssignRoles(menuID uint, roleIDs []uint) error
	RemoveRoles(menuID uint) error
	FindByPath(path string) (*model.Menu, error)
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) Delete(id uint) error {
	return r.db.Delete(&model.Menu{}, id).Error
}

func (r *menuRepository) FindByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.Preload("Roles").First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) FindAll(filter *model.MenuFilter) ([]model.Menu, int64, error) {
	var menus []model.Menu
	var total int64

	query := r.db.Model(&model.Menu{}).Preload("Roles")

	// Apply filters
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Path != "" {
		query = query.Where("path ILIKE ?", "%"+filter.Path+"%")
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// Count total records
	query.Count(&total)

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset)
	}

	// Order by order_index
	query = query.Order("order_index ASC, id ASC")

	err := query.Find(&menus).Error
	if err != nil {
		return nil, 0, err
	}

	return menus, total, nil
}

func (r *menuRepository) FindByUserRole(roleID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.
		Joins("JOIN menu_roles ON menu_roles.menu_id = menus.id").
		Where("menu_roles.role_id = ? AND menus.is_active = ?", roleID, true).
		Order("menus.order_index ASC, menus.id ASC").
		Preload("Children", "is_active = ?", true).
		Find(&menus).Error
	
	if err != nil {
		return nil, err
	}

	// Filter out menus that have parent_id (only return top-level menus)
	var topLevelMenus []model.Menu
	for _, menu := range menus {
		if menu.ParentID == nil {
			topLevelMenus = append(topLevelMenus, menu)
		}
	}

	return topLevelMenus, nil
}

func (r *menuRepository) AssignRoles(menuID uint, roleIDs []uint) error {
	// Start a transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Remove existing role assignments
	if err := tx.Where("menu_id = ?", menuID).Delete(&model.MenuRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create new role assignments
	for _, roleID := range roleIDs {
		menuRole := model.MenuRole{
			MenuID: menuID,
			RoleID: roleID,
		}
		if err := tx.Create(&menuRole).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *menuRepository) RemoveRoles(menuID uint) error {
	return r.db.Where("menu_id = ?", menuID).Delete(&model.MenuRole{}).Error
}

func (r *menuRepository) FindByPath(path string) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.Where("path = ?", path).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}
