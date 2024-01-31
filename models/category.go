package models

import "time"

type Category struct {
	ID             int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz         string     `gorm:"type:varchar(250) not null" json:"nameUz"`
	NameRu         string     `gorm:"type:varchar(250) not null" json:"nameRu"`
	NameEn         string     `gorm:"type:varchar(250) not null" json:"nameEn"`
	DescriptionRu  string     `gorm:"type:text not null" json:"descriptionRu"`
	DescriptionUz  string     `gorm:"type:text not null" json:"descriptionUz"`
	DescriptionEn  string     `gorm:"type:text not null" json:"descriptionEn"`
	Category       *Category  `gorm:"foreignKey:CategoryID" json:"parent"`
	CategoryID     *int       `gorm:"type:integer;default:null" json:"parentId"`
	Position       *int       `gorm:"type:integer;default:null" json:"position"`
	SeoTitle       string     `gorm:"type:varchar(150)" json:"seoTitle"`
	SeoDescription string     `gorm:"type:varchar(300)" json:"seoDescription"`
	Image          string     `gorm:"type:varchar(300)" json:"image"`
	IsActive       *bool      `gorm:"type:boolean;default:true;index" json:"isActive"`
	Created        *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID      *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt      *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
	Updated        *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID      *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt      *time.Time `gorm:"type:timestamptz;default:null" json:"updatedAt"`
	DeletedID      *int       `gorm:"type:bigint;default:null"  json:"-"`
	Deleted        *Admins    `gorm:"foreignKey:DeletedID"       json:"deleted"`
	DeletedAt      *time.Time `gorm:"type:timestamptz;default:null" json:"deletedAt"`
}

type ProductAdditions struct {
	ProductCategoryID  int       `gorm:"type:integer not null" json:"productCategoryID"`
	ProductCategory    *Category `gorm:"foreignKey:ProductCategoryID;unique_index:idx_index_addition" json:"-"`
	AdditionCategoryID int       `gorm:"type:integer not null" json:"additionCategoryID"`
	AdditionCategory   *Category `gorm:"foreignKey:AdditionCategoryID;unique_index:idx_index_addition" json:"-"`
}
type ProductAdditionRequest struct {
	ProductCategoryID  int `json:"productCategoryId" form:"productCategoryId"`
	AdditionCategoryID int `json:"additionCategoryId" form:"additionCategoryId"`
}

type CategoryRequest struct {
	NameRu         string `json:"nameRu" form:"nameRu"`
	NameEn         string `json:"nameEn" form:"nameEn"`
	NameUz         string `json:"nameUz" form:"nameUz"`
	Position       *int   `json:"position" form:"position"`
	DescriptionEn  string `json:"descriptionEn" form:"descriptionEn"`
	DescriptionRu  string `json:"descriptionRu" form:"descriptionRu"`
	DescriptionUz  string `json:"descriptionUz" form:"descriptionUz"`
	ParentID       int    `json:"parentId" form:"parentId"`
	IsActive       *bool  `json:"isActive" form:"isActive"`
	SeoTitle       string `json:"seoTitle" form:"seoTitle"`
	SeoDescription string `json:"seoDescription" form:"seoDescription"`
	Image          string `json:"-" form:"-"`
}
type CategoryFilter struct {
	Name     string `json:"name" form:"name"`
	ParentID *int   `json:"parentId" form:"parentId"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}
