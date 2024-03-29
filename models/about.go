package models

import "time"

type About struct {
	ID            int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz        string     `gorm:"type:varchar(250);default:null" json:"name_uz"`
	NameRu        string     `gorm:"type:varchar(250);default:null" json:"name_ru"`
	NameEn        string     `gorm:"type:varchar(250);default:null" json:"name_en"`
	DescriptionUz string     `gorm:"type:text;default:null" json:"description_uz"`
	DescriptionRu string     `gorm:"type:text;default:null" json:"description_ru"`
	DescriptionEn string     `gorm:"type:text;default:null" json:"description_en"`
	Position      *int       `gorm:"type:integer;default:null;index" json:"position"`
	Type          string     `gorm:"type:varchar(50) not null" json:"type"`
	Image         string     `gorm:"type:varchar(300);default:null" json:"image"`
	Created       *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID     *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt     *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
	Updated       *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID     *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt     *time.Time `gorm:"type:timestamptz;default:null" json:"updated_at"`
}

type AboutRequest struct {
	NameUz        string `json:"name_uz" form:"name_uz"`
	NameRu        string `json:"name_ru" form:"name_ru"`
	NameEn        string `json:"name_en" form:"name_en"`
	DescriptionEn string `json:"description_en" form:"description_en"`
	DescriptionRu string `json:"description_ru" form:"description_ru"`
	DescriptionUz string `json:"description_uz" form:"description_uz"`
	Position      *int   `json:"position" form:"position"`
	Type          string `json:"type" form:"type"`
	Image         string `json:"-" form:"-"`
}
