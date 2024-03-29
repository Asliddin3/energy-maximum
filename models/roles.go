package models

import "time"

type Roles struct {
	ID        int        `gorm:"type:bigint;primaryKey" json:"id"`
	Key       string     `gorm:"type:varchar(255) not null;unique"               json:"key"`
	Title     string     `gorm:"type:varchar(255) not null"                      json:"title"`
	Comment   string     `gorm:"type:text;default:null"                          json:"comment"`
	IsActive  bool       `gorm:"type:boolean;default:true;index"                json:"is_active"`
	IsDeleted bool       `gorm:"type:boolean;default:false;index"               json:"is_deleted"`
	CreatedID *int       `gorm:"type:bigint;default:null;index"  json:"-"`
	Created   *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null;index" json:"created_at"`
	Updated   *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"updated_at"`
}
type RoleItems struct {
	ID            int          `gorm:"type:bigint;primaryKey" json:"id"`
	Role          *Roles       `gorm:"foreignKey:RoleID" json:"role"`
	RoleID        int          `gorm:"type:integer;default:null;index" json:"role_id"`
	ModuleItemKey string       `gorm:"type:varchar(255); not null;index" json:"key"`
	ModuleItem    *ModuleItems `gorm:"foreignKey:ModuleItemKey;onDelete:CASCADE" json:"module_item"`
	Created       *Admins      `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID     *int         `gorm:"type:bigint;default:null;index"  json:"-"`
	CreatedAt     *time.Time   `gorm:"type:timestamptz;default:null;index" json:"created_at"`
	Updated       *Admins      `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID     *int         `gorm:"type:bigint;default:null;index"  json:"-"`
	UpdatedAt     *time.Time   `gorm:"type:timestamptz;default:null" json:"updated_at"`
}

type RoleItemRequest struct {
	RoleID int    `json:"role_id" form:"role_id"`
	Key    string `json:"key" form:"key"`
}

type RolesRequest struct {
	Title    string `json:"title" form:"title"`
	Key      string `json:"key" form:"key"`
	Comment  string `json:"comment" form:"comment"`
	IsActive bool   `json:"is_active" form:"is_active"`
}

type RoleFilter struct {
	Name     string `json:"name" form:"name"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	// All        bool   `json:"all"  form:"all"`
	WithDelete bool `json:"with_delete" form:"with_delete"`
}
type RoleItemFilter struct {
	Name     string `json:"name" form:"name"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	// All        bool   `json:"all"  form:"all"`
	WithDelete bool `json:"with_delete" form:"with_delete"`
}

type UpdateRoleItemsList struct {
	ModuleItemKeys []string `json:"module_item_keys" `
	RoleID         int      `json:"role_id"`
}
type RoleWithModuleItems struct {
	Roles
	ModuleItemKeys []string `json:"module_item_keys"`
}
