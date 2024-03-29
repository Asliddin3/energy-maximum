package models

import "time"

type Admins struct {
	ID        int        `gorm:"type:bigint;primaryKey" json:"id"`
	Username  string     `gorm:"type:varchar(255);unique" json:"username"`
	Password  string     `gorm:"type:varchar(255)" json:"-"`
	Role      *Roles     `gorm:"foreignKey:RoleID" json:"role"`
	RoleID    int        `gorm:"type:bigint;default:null" json:"role_id"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
	UpdatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"updated_at"`
	LastVisit *time.Time `gorm:"type:timestamptz;default:null" json:"last_visit"`
	IsActive  *bool      `gorm:"type:boolean;default:true;index" json:"is_active"`
	DeletedAt *time.Time `gorm:"type:timestamptz;default:null" json:"deleted_at"`
}
type AdminResponse struct {
	Admins
	ModuleItemKeys []string `json:"moduleItemKeys"`
}
type AdminMetadata struct {
	Id int
}
type AdminsCreateRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	RoleID   int    `json:"role_id" form:"role_id"`
	IsActive *bool  `json:"is_active" form:"is_active"`
}
type AdminsRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type AdminAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccessToken struct {
	AccessToken string `json:"accessToken"`
}

type TokenResponse struct {
	AccessToken    string   `json:"accessToken"`
	ModuleItemKeys []string `json:"moduleItemKeys"`
}
