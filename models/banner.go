package models

import "time"

type Banner struct {
	ID            int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz        string     `gorm:"type:varchar(250) not null" json:"nameUz"`
	NameRu        string     `gorm:"type:varchar(250) not null" json:"nameRu"`
	NameEn        string     `gorm:"type:varchar(250) not null" json:"nameEn"`
	DescriptionRu string     `gorm:"type:text not null" json:"descriptionRu"`
	DescriptionUz string     `gorm:"type:text not null" json:"descriptionUz"`
	DescriptionEn string     `gorm:"type:text not null" json:"descriptionEn"`
	Position      *int       `gorm:"type:integer;default:null" json:"position"`
	Image         string     `gorm:"type:varchar(300)" json:"image"`
	Created       *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID     *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt     *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
	Updated       *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID     *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt     *time.Time `gorm:"type:timestamptz;default:null" json:"updatedAt"`
}

type BannerRequest struct {
	NameRu        string `json:"nameRu" form:"nameRu"`
	NameEn        string `json:"nameEn" form:"nameEn"`
	NameUz        string `json:"nameUz" form:"nameUz"`
	Position      *int   `json:"position" form:"position"`
	DescriptionEn string `json:"descriptionEn" form:"descriptionEn"`
	DescriptionRu string `json:"descriptionRu" form:"descriptionRu"`
	DescriptionUz string `json:"descriptionUz" form:"descriptionUz"`
	Image         string `json:"-" form:"-"`
}
