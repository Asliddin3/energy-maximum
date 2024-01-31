package models

import "time"

type Products struct {
	ID             int        `gorm:"type:bigint not null;primaryKey" json:"id"`
	NameUz         string     `gorm:"type:varchar(250) not null" json:"nameUz"`
	NameRu         string     `gorm:"type:varchar(250) not null" json:"nameRu"`
	NameEn         string     `gorm:"type:varchar(250) not null" json:"nameEn"`
	DescriptionRu  string     `gorm:"type:text not null" json:"descriptionRu"`
	DescriptionUz  string     `gorm:"type:text not null" json:"descriptionUz"`
	DescriptionEn  string     `gorm:"type:text not null" json:"descriptionEn"`
	SeoTitle       string     `gorm:"type:varchar(150) not null" json:"seoTitle"`
	SeoDescription string     `gorm:"type:varchar(300) not null" json:"seoDescription"`
	Price          float64    `gorm:"type:decimal(16,2) not null" json:"price"`
	Parent         *Category  `gorm:"foreignKey:ParentID" json:"parent"`
	ParentID       *int       `gorm:"type:bigint;default:null" json:"parentId"`
	IsTop          *bool      `gorm:"type:boolean;default:null" json:"isTop"`
	IsNew          *bool      `gorm:"type:boolean;default:null" json:"isNew"`
	Position       *int       `gorm:"type:integer;default:null" json:"position"`
	Country        *Country   `gorm:"foreignKey:CountryID" json:"country"`
	CountryID      *int       `gorm:"type:bigint;default:null" json:"countryId"`
	Brand          *Brand     `gorm:"foreignKey:BrandID" json:"brand"`
	BrandID        *int       `gorm:"type:bigint;default:null" json:"brandId"`
	IsActive       *bool      `gorm:"type:boolean not null;index" json:"isActive"`
	Image          *string    `gorm:"type:varchar(300);default:null" json:"image"`
	Created        *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID      *int       `gorm:"type:integer;default:null"  json:"-"`
	CreatedAt      *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
	Updated        *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID      *int       `gorm:"type:integer;default:null"  json:"-"`
	UpdatedAt      *time.Time `gorm:"type:timestamptz;default:null" json:"updatedAt"`
	DeletedID      *int       `gorm:"type:integer;default:null"  json:"-"`
	Deleted        *Admins    `gorm:"foreignKey:DeletedID"       json:"deleted"`
	DeletedAt      *time.Time `gorm:"type:timestamptz;default:null" json:"deletedAt"`
}
type Parameters struct {
	ID        int        `gorm:"type:bigint not null;primaryKey" json:"id"`
	NameUz    string     `gorm:"type:varchar(250) not null" json:"nameUz"`
	NameRu    string     `gorm:"type:varchar(250) not null" json:"nameRu"`
	NameEn    string     `gorm:"type:varchar(250) not null" json:"nameEn"`
	Position  *int       `gorm:"type:integer;default:null" json:"position"`
	Created   *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID *int       `gorm:"type:integer;default:null"  json:"-"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
}

type ProductsIds struct {
	ProductsIds []int `json:"productsIds"`
}

type ParameterFilter struct {
	Name     string `json:"name" form:"name" `
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}

type ParameterResponse struct {
	Parameters []Parameters `json:"parameters" form:"parameters"`
	Page       int          `json:"page" form:"page"`
	PageSize   int          `json:"pageSize" form:"pageSize"`
}

type ParametersRequest struct {
	NameRu   string `json:"nameRu" form:"nameRu"`
	NameEn   string `json:"nameEn" form:"nameEn"`
	NameUz   string `json:"nameUz" form:"nameUz"`
	Position *int   `json:"position" form:"position"`
}

type ProductParameters struct {
	Product     *Products   `gorm:"foreignKey:ProductID" json:"product"`
	ProductID   int         `gorm:"type:bigint not null" json:"productId"`
	Parameter   *Parameters `gorm:"foreignKey:ParameterID" json:"parameter"`
	ParameterID int         `gorm:"type:bigint not null" json:"parameterId"`
	ValRu       string      `gorm:"type:varchar(255) not null" json:"valRu"`
	ValUz       string      `gorm:"type:varchar(255) not null" json:"valUz"`
	ValEn       string      `gorm:"type:varchar(255) not null" json:"valEn"`
}

type Recommend struct {
	ID        int        `gorm:"type:bigint not null;primaryKey" json:"id"`
	NameUz    string     `gorm:"type:varchar(250) not null" json:"nameUz"`
	NameRu    string     `gorm:"type:varchar(250) not null" json:"nameRu"`
	NameEn    string     `gorm:"type:varchar(250) not null" json:"nameEn"`
	Created   *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID *int       `gorm:"type:integer;default:null"  json:"-"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
}
type RecommendItems struct {
	Recommend   *Recommend `gorm:"foreignKey:RecommendID" json:"-"`
	RecommendID int        `gorm:"type:integer not null" json:"recommendId"`
	Product     *Products  `gorm:"type:foreignKey:ProductID" json:"-"`
	ProductID   int        `gorm:"type:integer not null" json:"productId"`
}

type RecommendRequest struct {
	NameRu   string `json:"nameRu" form:"nameRu"`
	NameEn   string `json:"nameEn" form:"nameEn"`
	NameUz   string `json:"nameUz" form:"nameUz"`
	IsAnalog bool   `json:"isAnalog" form:"isAnalog"`
}

type ProductRecommendResponse struct {
	*Products
	RecommendID int `json:"productRecommendId"`
}
type ProductRecommendRequest struct {
	ProductID   int  `json:"productId"`
	IsAddition  bool `json:"isAddition"`
	RecommendID int  `json:"recommendId"`
}
type ProductResponse struct {
	*Products
	Media      []ProductMedia      `json:"media"`
	Parameters []ProductParameters `json:"parameters"`
}

type ProductRequest struct {
	NameRu         string  `json:"nameRu" form:"nameRu"`
	NameEn         string  `json:"nameEn" form:"nameEn"`
	NameUz         string  `json:"nameUz" form:"nameUz"`
	DescriptionEn  string  `json:"descriptionEn" form:"descriptionEn"`
	DescriptionRu  string  `json:"descriptionRu" form:"descriptionRu"`
	DescriptionUz  string  `json:"descriptionUz" form:"descriptionUz"`
	IsTop          *bool   `json:"isTop" form:"isTop"`
	IsNew          *bool   `json:"isNew" form:"isNew"`
	Position       *int    `json:"position" form:"position"`
	Price          float64 `json:"price" form:"price"`
	ParentID       int     `json:"parentId" form:"parentId"`
	BrandID        int     `json:"brandId" form:"brandId"`
	CountryID      *int    `json:"countryId" form:"countryId"`
	SeoTitle       string  `json:"seoTitle" form:"seoTitle"`
	SeoDescription string  `json:"seoDescription" form:"seoDescription"`
	IsActive       *bool   `json:"isActive" form:"isActive"`
	Image          *string `json:"-" form:"- "`
}
type ProductParamReq struct {
	Parameters []ProductParametersRequest `json:"parameters" form:"parameters"`
}

type ProductParametersRequest struct {
	ParameterID int    `json:"parameterId" form:"parameterId"`
	ValRu       string `json:"valRu" form:"valRu"`
	ValEn       string `json:"valEn" form:"valEn"`
	ValUz       string `json:"valUz" form:"valUz"`
}

type ProductsFilter struct {
	ParentID    int     `json:"parentId" form:"parentId"`
	PriceFrom   float64 `json:"priceFrom" form:"priceFrom"`
	PriceTo     float64 `json:"priceTo" form:"priceTo"`
	IsTop       *bool   `json:"isTop" form:"isTop"`
	IsNew       *bool   `json:"isNew" form:"isNew"`
	MultiSearch string  `json:"multiSearch" form:"multiSearch"`
	Page        int     `json:"page" form:"page"`
	PageSize    int     `json:"pageSize" form:"pageSize"`
}

type ProductsList struct {
	Products []Products `json:"products"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
	Count    int        `json:"count"`
}
type ProductIds struct {
	ProductsIds []int `json:"productsIds" form:"productsIds"`
}

type ProductMedia struct {
	ID        int       `gorm:"type:bigint not null;primaryKey" json:"id"`
	Product   *Products `gorm:"foreignKey:ProductID" json:"-"`
	ProductID *int      `gorm:"type:bigint;default:" json:"productId"`
	Position  *int      `gorm:"type:integer;default:" json:"position"`
	Type      string    `gorm:"type:varchar(100)" json:"type"`
	Media     string    `gorm:"type:varchar(300)" json:"media"`
}

type ProductMediaRequest struct {
	Position int    `json:"position" form:"position" binding:"required"`
	Type     string `json:"type" form:"type" binding:"required"`
	Media    string `json:"-" form:"-"`
}
