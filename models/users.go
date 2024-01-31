package models

import "time"

type Admins struct {
	ID          int        `gorm:"type:bigint;primaryKey" json:"id"`
	Username    string     `gorm:"type:varchar(255);unique" json:"username"`
	Password    string     `gorm:"type:varchar(255)" json:"-"`
	IsSuperuser *bool      `gorm:"type:boolean not null" json:"isSuperuser"`
	CreatedAt   *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
	UpdatedAt   *time.Time `gorm:"type:timestamptz;default:null" json:"updatedAt"`
	LastVisit   *time.Time `gorm:"type:timestamptz;default:null" json:"lastVisit"`
	IsActive    *bool      `gorm:"type:boolean;default:true" json:"isActive"`
	DeletedAt   *time.Time `gorm:"type:timestamptz;default:null" json:"deletedAt"`
}
type AdminMetadata struct {
	Id          int
	IsSuperuser bool
}
type AdminsCreateRequest struct {
	Username    string `json:"username" form:"username"`
	Password    string `json:"password" form:"password"`
	IsSuperuser bool   `json:"isSuperuser" form:"isSuperuser"`
	IsActive    *bool  `json:"isActive" form:"isActive"`
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
	AccessToken string `json:"accessToken"`
}
