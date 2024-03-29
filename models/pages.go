package models

import "time"

type Pages struct {
	ID               int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz           string     `gorm:"type:varchar(250) not null" json:"name_uz"`
	NameRu           string     `gorm:"type:varchar(250) not null" json:"name_ru"`
	NameEn           string     `gorm:"type:varchar(250) not null" json:"name_en"`
	Url              string     `gorm:"type:varchar(250);default:null" json:"url"`
	DescriptionRu    string     `gorm:"type:text not null" json:"description_ru"`
	DescriptionUz    string     `gorm:"type:text not null" json:"description_uz"`
	DescriptionEn    string     `gorm:"type:text not null" json:"description_en"`
	Position         *int       `gorm:"type:integer;default:null;index" json:"position"`
	SeoTitleRu       string     `gorm:"type:varchar(150)" json:"seo_title_ru"`
	SeoTitleUz       string     `gorm:"type:varchar(150)" json:"seo_title_uz"`
	SeoTitleEn       string     `gorm:"type:varchar(150)" json:"seo_title_en"`
	SeoDescriptionRu string     `gorm:"type:varchar(300)" json:"seo_description_ru"`
	SeoDescriptionEn string     `gorm:"type:varchar(300)" json:"seo_description_en"`
	SeoDescriptionUz string     `gorm:"type:varchar(300)" json:"seo_description_uz"`
	Image            string     `gorm:"type:varchar(300)" json:"image"`
	IsActive         *bool      `gorm:"type:boolean;default:true;index" json:"is_active"`
	Created          *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID        *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt        *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
	Updated          *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID        *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt        *time.Time `gorm:"type:timestamptz;default:null" json:"updated_at"`
	DeletedID        *int       `gorm:"type:bigint;default:null"  json:"-"`
	Deleted          *Admins    `gorm:"foreignKey:DeletedID"       json:"deleted"`
	DeletedAt        *time.Time `gorm:"type:timestamptz;default:null" json:"deleted_at"`
}

type PagesRequest struct {
	NameRu           string `json:"name_ru" form:"name_ru"`
	NameEn           string `json:"name_en" form:"name_en"`
	NameUz           string `json:"name_uz" form:"name_uz"`
	Position         *int   `json:"position" form:"position"`
	Url              string `json:"url" form:"url"`
	DescriptionEn    string `json:"description_en" form:"description_en"`
	DescriptionRu    string `json:"description_ru" form:"description_ru"`
	DescriptionUz    string `json:"description_uz" form:"description_uz"`
	IsActive         *bool  `json:"is_active" form:"is_active"`
	SeoTitleRu       string `json:"seo_title_ru" form:"seo_title_ru"`
	SeoTitleUz       string `json:"seo_title_uz" form:"seo_title_uz"`
	SeoTitleEn       string `json:"seo_title_en" form:"seo_title_en"`
	SeoDescriptionRu string `json:"seo_description_ru" form:"seo_description_ru"`
	SeoDescriptionUz string `json:"seo_description_uz" form:"seo_description_uz"`
	SeoDescriptionEn string `json:"seo_description_en" form:"seo_description_en"`
	Image            string `json:"-" form:"-"`
}
type PagesFilter struct {
	Name     string `json:"name" form:"name"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
}
