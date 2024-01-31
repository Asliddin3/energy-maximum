package models

import "time"

type Contact struct {
	ID           int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz       string     `gorm:"type:varchar(250) not null" json:"nameUz"`
	NameRu       string     `gorm:"type:varchar(250) not null" json:"nameRu"`
	NameEn       string     `gorm:"type:varchar(250) not null" json:"nameEn"`
	IsMain       *bool      `gorm:"type:boolean not null" json:"isMain"`
	Phone        string     `gorm:"type:varchar(250) not null" json:"phone"`
	Address      string     `gorm:"type:varchar(250) not null" json:"address"`
	Longitude    string     `gorm:"type:varchar(250) not null" json:"longitude"`
	Latitude     string     `gorm:"type:varchar(250) not null" json:"latitude"`
	Email        string     `gorm:"type:varchar(250) not null" json:"email"`
	WorkingHours string     `gorm:"type:varchar(250) not null" json:"workingHours"`
	Created      *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID    *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt    *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
}

type ContactRequest struct {
	NameUz       string `json:"nameUz" form:"nameUz"`
	NameRu       string `json:"nameRu" form:"nameRu"`
	NameEn       string `json:"nameEn" form:"nameEn"`
	IsMain       *bool  `json:"isMain" form:"isMain"`
	Phone        string `json:"phone" form:"phone"`
	Address      string `json:"address" form:"address"`
	Longitude    string `json:"longitude" form:"longitude"`
	Latitude     string `json:"latitude" form:"latitude"`
	Email        string `json:"email" form:"email"`
	WorkingHours string `json:"workingHours" form:"workingHours"`
}
