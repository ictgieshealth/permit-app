package model

import "time"

type User struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID      int64     `gorm:"not null" json:"role_id"`
	Username    string    `gorm:"type:varchar(100);not null;unique" json:"username"`
	Email       string    `gorm:"type:varchar(255);not null;unique" json:"email"`
	Password    string    `gorm:"type:varchar(255);not null" json:"-"`
	FullName    string    `gorm:"type:varchar(255);not null" json:"full_name"`
	PhoneNumber string    `gorm:"type:varchar(20)" json:"phone_number"`
	NIP         string    `gorm:"type:varchar(50)" json:"nip"`
	IsActive    bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Role        *Role         `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID"`
	UserDomains []UserDomain  `json:"user_domains,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Domains     []*Domain     `json:"domains,omitempty" gorm:"many2many:user_domains;"`
}

type UserDomain struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null" json:"user_id"`
	DomainID  int64     `gorm:"not null" json:"domain_id"`
	IsDefault bool      `gorm:"type:boolean;not null;default:false" json:"is_default"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Domain *Domain `json:"domain,omitempty" gorm:"foreignKey:DomainID;references:ID"`
}

type UserRequest struct {
	RoleID      int64   `json:"role_id" validate:"required"`
	Username    string  `json:"username" validate:"required,max=100"`
	Email       string  `json:"email" validate:"required,email,max=255"`
	Password    string  `json:"password" validate:"required,min=6"`
	FullName    string  `json:"full_name" validate:"required,max=255"`
	PhoneNumber string  `json:"phone_number" validate:"omitempty,max=20"`
	NIP         string  `json:"nip" validate:"omitempty,max=50"`
	IsActive    *bool   `json:"is_active"`
	DomainIDs   []int64 `json:"domain_ids" validate:"required,min=1"`
}

type UserUpdateRequest struct {
	RoleID      int64   `json:"role_id" validate:"omitempty"`
	Username    string  `json:"username" validate:"omitempty,max=100"`
	Email       string  `json:"email" validate:"omitempty,email,max=255"`
	Password    string  `json:"password" validate:"omitempty,min=6"`
	FullName    string  `json:"full_name" validate:"omitempty,max=255"`
	PhoneNumber string  `json:"phone_number" validate:"omitempty,max=20"`
	NIP         string  `json:"nip" validate:"omitempty,max=50"`
	IsActive    *bool   `json:"is_active"`
	DomainIDs   []int64 `json:"domain_ids" validate:"omitempty,min=1"`
}

type UserResponse struct {
	ID          int64            `json:"id"`
	RoleID      int64            `json:"role_id"`
	Username    string           `json:"username"`
	Email       string           `json:"email"`
	FullName    string           `json:"full_name"`
	PhoneNumber string           `json:"phone_number"`
	NIP         string           `json:"nip"`
	IsActive    bool             `json:"is_active"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Role        *RoleResponse    `json:"role,omitempty"`
	Domains     []DomainResponse `json:"domains,omitempty"`
}

type UserListRequest struct {
	DomainID *int64 `form:"domain_id"`
	RoleID   *int64 `form:"role_id"`
	Username string `form:"username"`
	Email    string `form:"email"`
	FullName string `form:"full_name"`
	IsActive *bool  `form:"is_active"`
	Page     int    `form:"page" validate:"omitempty,min=1"`
	Limit    int    `form:"limit" validate:"omitempty,min=1,max=10000"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	DomainID *int64 `json:"domain_id"` // Optional: user can specify which domain to login to
}

type LoginResponse struct {
	Token        string        `json:"token"`
	User         *UserResponse `json:"user"`
	DefaultDomain *DomainResponse `json:"default_domain,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UserDomainRequest struct {
	DomainID  int64 `json:"domain_id" validate:"required"`
	IsDefault bool  `json:"is_default"`
}

type UpdateProfileRequest struct {
	FullName    string `json:"full_name" validate:"omitempty,max=255"`
	Email       string `json:"email" validate:"omitempty,email,max=255"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,max=20"`
	NIP         string `json:"nip" validate:"omitempty,max=50"`
}

