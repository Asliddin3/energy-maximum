package models

import "time"

type Contact struct {
	ID           int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz       string     `gorm:"type:varchar(250) not null" json:"name_uz"`
	NameRu       string     `gorm:"type:varchar(250) not null" json:"name_ru"`
	NameEn       string     `gorm:"type:varchar(250) not null" json:"name_en"`
	IsMain       *bool      `gorm:"type:boolean not null" json:"is_main"`
	Phone        string     `gorm:"type:varchar(250) not null" json:"phone"`
	Address      string     `gorm:"type:varchar(250) not null" json:"address"`
	Longitude    string     `gorm:"type:varchar(250) not null" json:"longitude"`
	Latitude     string     `gorm:"type:varchar(250) not null" json:"latitude"`
	Email        string     `gorm:"type:varchar(250) not null" json:"email"`
	WorkingHours string     `gorm:"type:varchar(250) not null" json:"working_hours"`
	Created      *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID    *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt    *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
}

type ContactRequest struct {
	NameUz       string `json:"name_uz" form:"name_uz"`
	NameRu       string `json:"name_ru" form:"name_ru"`
	NameEn       string `json:"name_en" form:"name_en"`
	IsMain       *bool  `json:"is_main" form:"is_main"`
	Phone        string `json:"phone" form:"phone"`
	Address      string `json:"address" form:"address"`
	Longitude    string `json:"longitude" form:"longitude"`
	Latitude     string `json:"latitude" form:"latitude"`
	Email        string `json:"email" form:"email"`
	WorkingHours string `json:"working_hours" form:"working_hours"`
}
