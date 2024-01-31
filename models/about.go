package models

import "time"

type About struct {
	ID            int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz        string     `gorm:"type:varchar(250);default:null" json:"nameUz"`
	NameRu        string     `gorm:"type:varchar(250);default:null" json:"nameRu"`
	NameEn        string     `gorm:"type:varchar(250);default:null" json:"nameEn"`
	DescriptionUz string     `gorm:"type:text;default:null" json:"descriptionUz"`
	DescriptionRu string     `gorm:"type:text;default:null" json:"descriptionRu"`
	DescriptionEn string     `gorm:"type:text;default:null" json:"descriptionEn"`
	Position      *int       `gorm:"type:integer;default:null" json:"position"`
	Type          string     `gorm:"type:varchar(50) not null" json:"type"`
	Image         string     `gorm:"type:varchar(300);default:null" json:"image"`
	Created       *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID     *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt     *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
	Updated       *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID     *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt     *time.Time `gorm:"type:timestamptz;default:null" json:"updatedAt"`
}

type AboutRequest struct {
	NameUz        string `json:"nameUz" form:"nameUz"`
	NameRu        string `json:"nameRu" form:"nameRu"`
	NameEn        string `json:"nameEn" form:"nameEn"`
	DescriptionEn string `json:"descriptionEn" form:"descriptionEn"`
	DescriptionRu string `json:"descriptionRu" form:"descriptionRu"`
	DescriptionUz string `json:"descriptionUz" form:"descriptionUz"`
	Position      *int   `json:"position" form:"position"`
	Type          string `json:"type" form:"type"`
	Image         string `json:"-" form:"-"`
}
