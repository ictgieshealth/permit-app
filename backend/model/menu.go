package model

import "time"

type Menu struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Path        string    `gorm:"size:255;uniqueIndex;not null" json:"path"`
	Icon        string    `gorm:"size:50" json:"icon,omitempty"`
	ParentID    *uint     `gorm:"index" json:"parent_id,omitempty"`
	OrderIndex  int       `gorm:"default:0;not null" json:"order_index"`
	IsActive    bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	Parent      *Menu     `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Menu    `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Roles       []Role    `gorm:"many2many:menu_roles" json:"roles,omitempty"`
}

type MenuRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MenuID    uint      `gorm:"not null;index" json:"menu_id"`
	RoleID    uint      `gorm:"not null;index" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	Menu      Menu      `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName specifies the table name for Menu model
func (Menu) TableName() string {
	return "menus"
}

// TableName specifies the table name for MenuRole model
func (MenuRole) TableName() string {
	return "menu_roles"
}

// MenuRequest represents the request payload for creating/updating a menu
type MenuRequest struct {
	Name       string  `json:"name" binding:"required"`
	Path       string  `json:"path" binding:"required"`
	Icon       string  `json:"icon"`
	ParentID   *uint   `json:"parent_id"`
	OrderIndex int     `json:"order_index"`
	RoleIDs    []uint  `json:"role_ids" binding:"required,min=1"`
}

// MenuFilter represents filter parameters for querying menus
type MenuFilter struct {
	Name     string `form:"name"`
	Path     string `form:"path"`
	IsActive *bool  `form:"is_active"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}
