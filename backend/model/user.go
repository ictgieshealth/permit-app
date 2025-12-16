package model

import "time"

type User struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string    `gorm:"type:varchar(100);not null;unique" json:"username"`
	Email       string    `gorm:"type:varchar(255);not null;unique" json:"email"`
	Password    string    `gorm:"type:varchar(255);not null" json:"-"`
	FullName    string    `gorm:"type:varchar(255);not null" json:"full_name"`
	PhoneNumber string    `gorm:"type:varchar(20)" json:"phone_number"`
	Nip         string    `gorm:"type:varchar(50)" json:"nip"`
	IsActive    bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	UserDomainRoles []UserDomainRole `json:"user_domain_roles,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

type UserDomainRole struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null" json:"user_id"`
	DomainID  int64     `gorm:"not null" json:"domain_id"`
	RoleID    int64     `gorm:"not null" json:"role_id"`
	IsDefault bool      `gorm:"type:boolean;not null;default:false" json:"is_default"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Domain *Domain `json:"domain,omitempty" gorm:"foreignKey:DomainID;references:ID"`
	Role   *Role   `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID"`
}

type UserRequest struct {
	Username    string                    `json:"username" validate:"required,max=100"`
	Email       string                    `json:"email" validate:"required,email,max=255"`
	Password    string                    `json:"password" validate:"required,min=6"`
	FullName    string                    `json:"full_name" validate:"required,max=255"`
	PhoneNumber string                    `json:"phone_number" validate:"omitempty,max=20"`
	Nip         string                    `json:"nip" validate:"omitempty,max=50"`
	IsActive    *bool                     `json:"is_active"`
	DomainRoles []UserDomainRoleRequest   `json:"domain_roles" validate:"required,min=1"`
}

type UserDomainRoleRequest struct {
	DomainID  int64 `json:"domain_id" validate:"required"`
	RoleID    int64 `json:"role_id" validate:"required"`
	IsDefault bool  `json:"is_default"`
}

type UserUpdateRequest struct {
	Username    string                    `json:"username" validate:"omitempty,max=100"`
	Email       string                    `json:"email" validate:"omitempty,email,max=255"`
	Password    string                    `json:"password" validate:"omitempty,min=6"`
	FullName    string                    `json:"full_name" validate:"omitempty,max=255"`
	PhoneNumber string                    `json:"phone_number" validate:"omitempty,max=20"`
	Nip         string                    `json:"nip" validate:"omitempty,max=50"`
	IsActive    *bool                     `json:"is_active"`
	DomainRoles []UserDomainRoleRequest   `json:"domain_roles" validate:"omitempty,min=1"`
}

type UserResponse struct {
	ID              int64                     `json:"id"`
	Username        string                    `json:"username"`
	Email           string                    `json:"email"`
	FullName        string                    `json:"full_name"`
	PhoneNumber     string                    `json:"phone_number"`
	Nip             string                    `json:"nip"`
	IsActive        bool                      `json:"is_active"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	DomainRoles     []UserDomainRoleResponse  `json:"domain_roles,omitempty"`
}

type UserDomainRoleResponse struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	DomainID  int64           `json:"domain_id"`
	RoleID    int64           `json:"role_id"`
	IsDefault bool            `json:"is_default"`
	Domain    *DomainResponse `json:"domain,omitempty"`
	Role      *RoleResponse   `json:"role,omitempty"`
}

type UserListRequest struct {
	DomainID *int64 `form:"domain_id"`
	RoleID   *int64 `form:"role_id"`
	Username string `form:"username"`
	Email    string `form:"email"`
	FullName string `form:"full_name"`
	IsActive *bool  `form:"is_active"`
	Category string `form:"category" validate:"omitempty,oneof=Permit Ticketing"` // Filter by role category
	Page     int    `form:"page" validate:"omitempty,min=1"`
	Limit    int    `form:"limit" validate:"omitempty,min=1,max=10000"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	DomainID *int64 `json:"domain_id"` // Optional: user can specify which domain to login to
}

type LoginResponse struct {
	Token         string                  `json:"token"`
	User          *UserResponse           `json:"user"`
	CurrentDomain *DomainResponse         `json:"current_domain"`
	CurrentRole   *RoleResponse           `json:"current_role"`
	Domains       []UserDomainRoleResponse `json:"domains"` // All domains user has access to
}

type SwitchDomainRequest struct {
	DomainID int64 `json:"domain_id" validate:"required"`
}

type SwitchDomainResponse struct {
	Token         string          `json:"token"`
	CurrentDomain *DomainResponse `json:"current_domain"`
	CurrentRole   *RoleResponse   `json:"current_role"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdateProfileRequest struct {
	FullName    string `json:"full_name" validate:"omitempty,max=255"`
	Email       string `json:"email" validate:"omitempty,email,max=255"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,max=20"`
	Nip         string `json:"nip" validate:"omitempty,max=50"`
}

