package models

import "time"

type Brand struct {
	ID             int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz         string     `gorm:"type:varchar(250) not null" json:"nameUz"`
	NameRu         string     `gorm:"type:varchar(250) not null" json:"nameRu"`
	NameEn         string     `gorm:"type:varchar(250) not null" json:"nameEn"`
	Letter         string     `gorm:"type:varchar(5);default:'';index" json:"letter"`
	DescriptionRu  string     `gorm:"type:text not null" json:"descriptionRu"`
	DescriptionUz  string     `gorm:"type:text not null" json:"descriptionUz"`
	DescriptionEn  string     `gorm:"type:text not null" json:"descriptionEn"`
	Position       *int       `gorm:"type:integer;default:null" json:"position"`
	SeoTitle       string     `gorm:"type:varchar(150)" json:"seoTitle"`
	SeoDescription string     `gorm:"type:varchar(300)" json:"seoDescription"`
	Image          string     `gorm:"type:varchar(300)" json:"image"`
	IsActive       *bool      `gorm:"type:boolean;default:true;index" json:"isActive"`
	Created        *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID      *int       `gorm:"type:bigint;default:null"    json:"-"`
	CreatedAt      *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
	Updated        *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID      *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt      *time.Time `gorm:"type:timestamptz;default:null" json:"updatedAt"`
}
type BrandRequest struct {
	NameRu         string `json:"nameRu" form:"nameRu"`
	NameEn         string `json:"nameEn" form:"nameEn"`
	NameUz         string `json:"nameUz" form:"nameUz"`
	Letter         string `json:"letter" form:"letter"`
	Position       *int   `json:"position" form:"position"`
	DescriptionEn  string `json:"descriptionEn" form:"descriptionEn"`
	DescriptionRu  string `json:"descriptionRu" form:"descriptionRu"`
	DescriptionUz  string `json:"descriptionUz" form:"descriptionUz"`
	SeoTitle       string `json:"seoTitle" form:"seoTitle"`
	SeoDescription string `json:"seoDescription" form:"seoDescription"`
	IsActive       *bool  `json:"isActive" form:"isActive"`
}

type BrandsByLetter struct {
	Letter string  `json:"letter" form:"letter"`
	Brands []Brand `json:"brands" form:"brands"`
}

type BrandFilter struct {
	Name     string `json:"name" form:"name"`
	Letter   string `json:"letter" form:"letter"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}

// type BrandFilterByLetter struct {
// }
