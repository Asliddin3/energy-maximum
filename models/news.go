package models

import "time"

type News struct {
	ID            int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz        string     `gorm:"type:varchar(250) not null;index" json:"name_uz"`
	NameRu        string     `gorm:"type:varchar(250) not null;index" json:"name_ru"`
	NameEn        string     `gorm:"type:varchar(250) not null;index" json:"name_en"`
	DescriptionRu string     `gorm:"type:text not null" json:"description_ru"`
	DescriptionUz string     `gorm:"type:text not null" json:"description_uz"`
	DescriptionEn string     `gorm:"type:text not null" json:"description_en"`
	Position      *int       `gorm:"type:integer;default:null;index" json:"position"`
	Image         string     `gorm:"type:varchar(300)" json:"image"`
	Created       *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID     *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt     *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
	Updated       *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID     *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt     *time.Time `gorm:"type:timestamptz;default:null" json:"updated_at"`
}

type NewsRequest struct {
	NameRu        string `json:"name_ru" form:"name_ru"`
	NameEn        string `json:"name_en" form:"name_en"`
	NameUz        string `json:"name_uz" form:"name_uz"`
	Position      *int   `json:"position" form:"position"`
	DescriptionEn string `json:"description_en" form:"description_en"`
	DescriptionRu string `json:"description_ru" form:"description_ru"`
	DescriptionUz string `json:"description_uz" form:"description_uz"`
	Image         string `json:"-" form:"-"`
}

type NewsFilter struct {
	Name     string `json:"name" form:"name"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
}
