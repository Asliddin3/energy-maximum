package models

import "time"

type Brand struct {
	ID             int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz         string     `gorm:"type:varchar(250) not null" json:"name_uz"`
	NameRu         string     `gorm:"type:varchar(250) not null" json:"name_ru"`
	NameEn         string     `gorm:"type:varchar(250) not null" json:"name_en"`
	Letter         string     `gorm:"type:varchar(5);default:'';index" json:"letter"`
	DescriptionRu  string     `gorm:"type:text not null" json:"description_ru"`
	DescriptionUz  string     `gorm:"type:text not null" json:"description_uz"`
	DescriptionEn  string     `gorm:"type:text not null" json:"description_en"`
	Position       *int       `gorm:"type:integer;default:null;index" json:"position"`
	SeoTitle       string     `gorm:"type:varchar(150)" json:"seo_title"`
	SeoDescription string     `gorm:"type:varchar(300)" json:"seo_description"`
	Image          string     `gorm:"type:varchar(300)" json:"image"`
	IsActive       *bool      `gorm:"type:boolean;default:true;index" json:"is_active"`
	Created        *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID      *int       `gorm:"type:bigint;default:null"    json:"-"`
	CreatedAt      *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
	Updated        *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID      *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt      *time.Time `gorm:"type:timestamptz;default:null" json:"updated_at"`
}
type BrandRequest struct {
	NameRu         string `json:"name_ru" form:"name_ru"`
	NameEn         string `json:"name_en" form:"name_en"`
	NameUz         string `json:"name_uz" form:"name_uz"`
	Letter         string `json:"letter" form:"letter"`
	Position       *int   `json:"position" form:"position"`
	DescriptionEn  string `json:"description_en" form:"description_en"`
	DescriptionRu  string `json:"description_ru" form:"description_ru"`
	DescriptionUz  string `json:"description_uz" form:"description_uz"`
	SeoTitle       string `json:"seo_title" form:"seo_title"`
	SeoDescription string `json:"seo_description" form:"seo_description"`
	IsActive       *bool  `json:"is_active" form:"is_active"`
}

type BrandsByLetter struct {
	Letter string  `json:"letter" form:"letter"`
	Brands []Brand `json:"brands" form:"brands"`
}

type BrandFilter struct {
	Name     string `json:"name" form:"name"`
	Letter   string `json:"letter" form:"letter"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
}

// type BrandFilterByLetter struct {
// }
