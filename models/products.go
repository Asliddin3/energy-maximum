package models

import "time"

type Products struct {
	ID               int        `gorm:"type:bigint not null;primaryKey" json:"id"`
	NameUz           string     `gorm:"type:varchar(250);default:null;index" json:"name_uz"`
	NameRu           string     `gorm:"type:varchar(250);default:null;index" json:"name_ru"`
	NameEn           string     `gorm:"type:varchar(250);default:null;index" json:"name_en"`
	Url              string     `gorm:"type:varchar(250);default:null;unique" json:"url"`
	DescriptionRu    string     `gorm:"type:text;default:null;index" json:"description_ru"`
	DescriptionUz    string     `gorm:"type:text;default:null" json:"description_uz"`
	DescriptionEn    string     `gorm:"type:text;default:null" json:"description_en"`
	SeoTitleRu       string     `gorm:"type:varchar(150);default:null" json:"seo_title_ru"`
	SeoTitleEn       string     `gorm:"type:varchar(150);default:null" json:"seo_title_en"`
	SeoTitleUz       string     `gorm:"type:varchar(150);default:null" json:"seo_title_uz"`
	SeoDescriptionRu string     `gorm:"type:varchar(300);default:null" json:"seo_description_ru"`
	SeoDescriptionEn string     `gorm:"type:varchar(300);default:null" json:"seo_description_en"`
	SeoDescriptionUz string     `gorm:"type:varchar(300);default:null" json:"seo_description_uz"`
	Price            float64    `gorm:"type:decimal(16,2) not null;index" json:"price"`
	Parent           *Category  `gorm:"foreignKey:ParentID" json:"parent"`
	ParentID         *int       `gorm:"type:bigint;default:null;index" json:"parent_id"`
	IsTop            *bool      `gorm:"type:boolean;default:false;index" json:"is_top"`
	IsNew            *bool      `gorm:"type:boolean;default:false;index" json:"is_new"`
	Position         *int       `gorm:"type:integer;default:null;index" json:"position"`
	Country          *Country   `gorm:"foreignKey:CountryID" json:"country"`
	CountryID        *int       `gorm:"type:bigint;default:null;index" json:"country_id"`
	Brand            *Brand     `gorm:"foreignKey:BrandID" json:"brand"`
	BrandID          *int       `gorm:"type:bigint;default:null;index" json:"brand_id"`
	IsActive         *bool      `gorm:"type:boolean not null;index" json:"is_active"`
	Image            string     `gorm:"type:varchar(300);default:null" json:"image"`
	Created          *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID        *int       `gorm:"type:integer;default:null"  json:"-"`
	CreatedAt        *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
	Updated          *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID        *int       `gorm:"type:integer;default:null"  json:"-"`
	UpdatedAt        *time.Time `gorm:"type:timestamptz;default:null" json:"updated_at"`
	DeletedID        *int       `gorm:"type:integer;default:null"  json:"-"`
	Deleted          *Admins    `gorm:"foreignKey:DeletedID"       json:"deleted"`
	DeletedAt        *time.Time `gorm:"type:timestamptz;default:null" json:"deleted_at"`
}
type Parameters struct {
	ID        int        `gorm:"type:bigint not null;primaryKey" json:"id"`
	NameUz    string     `gorm:"type:varchar(250) not null;index" json:"name_uz"`
	NameRu    string     `gorm:"type:varchar(250) not null;index" json:"name_ru"`
	NameEn    string     `gorm:"type:varchar(250) not null;index" json:"name_en"`
	Position  *int       `gorm:"type:integer;default:null;index" json:"position"`
	IsDeleted bool       `gorm:"type:boolean;default:false" json:"is_deleted"`
	Created   *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID *int       `gorm:"type:integer;default:null"  json:"-"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
}

type ProductsIds struct {
	ProductsIds []int `json:"productsIds"`
}

type ParameterFilter struct {
	Name        string `json:"name" form:"name" `
	Page        int    `json:"page" form:"page"`
	PageSize    int    `json:"page_size" form:"page_size"`
	WithDeleted bool   `json:"with_deleted" form:"with_deleted"`
}

type ParameterResponse struct {
	Parameters []Parameters `json:"parameters" form:"parameters"`
	Page       int          `json:"page" form:"page"`
	PageSize   int          `json:"page_size" form:"page_size"`
	Count      int          `json:"count" form:"count"`
}

type ParametersRequest struct {
	NameRu   string `json:"name_ru" form:"name_ru"`
	NameEn   string `json:"name_en" form:"name_en"`
	NameUz   string `json:"name_uz" form:"name_uz"`
	Position *int   `json:"position" form:"position"`
}

type ProductParameters struct {
	Product     *Products   `gorm:"foreignKey:ProductID" json:"product"`
	ProductID   int         `gorm:"type:bigint not null;index" json:"product_id"`
	Parameter   *Parameters `gorm:"foreignKey:ParameterID" json:"parameter"`
	ParameterID int         `gorm:"type:bigint not null"       json:"parameter_id"`
	ValRu       string      `gorm:"type:varchar(255) not null" json:"val_ru"`
	ValUz       string      `gorm:"type:varchar(255) not null" json:"val_uz"`
	ValEn       string      `gorm:"type:varchar(255) not null" json:"val_en"`
}

type ProductParameterResponse struct {
	ValRu       string `json:"val_ru"`
	ValUz       string `json:"val_uz"`
	ValEn       string `json:"val_en"`
	Position    int    `json:"position"`
	NameUz      string `json:"name_uz"`
	NameRu      string `json:"name_ru"`
	NameEn      string `json:"name_en"`
	ParameterID int    `json:"parameter_id"`
}

type Recommend struct {
	ID        int        `gorm:"type:bigint not null;primaryKey" json:"id"`
	NameUz    string     `gorm:"type:varchar(250) not null" json:"name_uz"`
	NameRu    string     `gorm:"type:varchar(250) not null" json:"name_ru"`
	NameEn    string     `gorm:"type:varchar(250) not null" json:"name_en"`
	Created   *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID *int       `gorm:"type:integer;default:null"  json:"-"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
}
type RecommendItems struct {
	Recommend   *Recommend `gorm:"foreignKey:RecommendID" json:"-"`
	RecommendID int        `gorm:"type:integer not null" json:"recommendId"`
	Product     *Products  `gorm:"type:foreignKey:ProductID" json:"-"`
	ProductID   int        `gorm:"type:integer not null" json:"product_id"`
}

type RecommendRequest struct {
	NameRu   string `json:"name_ru" form:"name_ru"`
	NameEn   string `json:"name_en" form:"name_en"`
	NameUz   string `json:"name_uz" form:"name_uz"`
	IsAnalog bool   `json:"isAnalog" form:"isAnalog"`
}

type ProductRecommendResponse struct {
	*Products
	RecommendID int `json:"productRecommendId"`
}
type ProductRecommendRequest struct {
	ProductID   int  `json:"product_id"`
	IsAddition  bool `json:"isAddition"`
	RecommendID int  `json:"recommendId"`
}
type ProductResponse struct {
	*Products
	Media      []ProductMedia             `json:"media"`
	Parameters []ProductParameterResponse `json:"parameters"`
}

type ProductRequest struct {
	NameRu           string  `json:"name_ru" form:"name_ru" binding:"required"`
	NameEn           string  `json:"name_en" form:"name_en"`
	NameUz           string  `json:"name_uz" form:"name_uz"`
	Url              string  `json:"url" form:"url"`
	DescriptionEn    string  `json:"description_en" form:"description_en"`
	DescriptionRu    string  `json:"description_ru" form:"description_ru"`
	DescriptionUz    string  `json:"description_uz" form:"description_uz"`
	IsTop            *bool   `json:"is_top" form:"is_top"`
	IsNew            *bool   `json:"is_new" form:"is_new"`
	Position         *int    `json:"position" form:"position"`
	Price            float64 `json:"price" form:"price"`
	ParentID         int     `json:"parent_id" form:"parent_id"`
	BrandID          int     `json:"brand_id" form:"brand_id"`
	CountryID        *int    `json:"country_id" form:"country_id"`
	SeoTitleRu       string  `json:"seo_title_ru" form:"seo_title_ru"`
	SeoTitleEn       string  `json:"seo_title_en" form:"seo_title_en"`
	SeoTitleUz       string  `json:"seo_title_uz" form:"seo_title_uz"`
	SeoDescriptionRu string  `json:"seo_description_ru" form:"seo_description_ru"`
	SeoDescriptionEn string  `json:"seo_description_en" form:"seo_description_en"`
	SeoDescriptionUz string  `json:"seo_description_uz" form:"seo_description_uz"`
	IsActive         *bool   `json:"is_active" form:"is_active"`
	Image            *string `json:"-" form:"- "`
}
type ProductParamReq struct {
	Parameters []ProductParametersRequest `json:"parameters" form:"parameters"`
}
type ProductParamDeleteReq struct {
	Parameters []int `json:"parameterIds" form:"parameterIds"`
}

type ProductParametersRequest struct {
	ParameterID int    `json:"parameterId" form:"parameterId"`
	ValRu       string `json:"valRu" form:"valRu"`
	ValEn       string `json:"valEn" form:"valEn"`
	ValUz       string `json:"valUz" form:"valUz"`
}

type ProductsFilter struct {
	ParentID    int     `json:"parent_id" form:"parent_id"`
	PriceFrom   float64 `json:"priceFrom" form:"priceFrom"`
	PriceTo     float64 `json:"priceTo" form:"priceTo"`
	IsTop       *bool   `json:"is_top" form:"is_top"`
	IsActive    *bool   `json:"is_active" form:"is_active"`
	IsNew       *bool   `json:"is_new" form:"is_new"`
	MultiSearch string  `json:"multiSearch" form:"multiSearch"`
	Page        int     `json:"page" form:"page"`
	PageSize    int     `json:"page_size" form:"page_size"`
}

type ProductsList struct {
	Products []Products `json:"products"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Count    int        `json:"count"`
}
type ProductIds struct {
	ProductsIds []int `json:"productsIds" form:"productsIds"`
}

type ProductMedia struct {
	ID        int       `gorm:"type:bigint not null;primaryKey" json:"id"`
	Product   *Products `gorm:"foreignKey:ProductID" json:"-"`
	ProductID *int      `gorm:"type:bigint;default:" json:"product_id"`
	Position  *int      `gorm:"type:integer;default:" json:"position"`
	Type      string    `gorm:"type:varchar(100)" json:"type"`
	Media     string    `gorm:"type:varchar(300)" json:"media"`
}

type ProductMediaRequest struct {
	Position int    `json:"position" form:"position" binding:"required"`
	Type     string `json:"type" form:"type" binding:"required"`
	Media    string `json:"-" form:"-"`
}
