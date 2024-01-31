package models

import "time"

type Vacancy struct {
	ID               int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz           string     `gorm:"type:varchar(250) not null" json:"nameUz"`
	NameRu           string     `gorm:"type:varchar(250) not null" json:"nameRu"`
	NameEn           string     `gorm:"type:varchar(250) not null" json:"nameEn"`
	DescriptionRu    string     `gorm:"type:text;default:null" json:"descriptionRu"`
	DescriptionUz    string     `gorm:"type:text;default:null" json:"descriptionUz"`
	DescriptionEn    string     `gorm:"type:text;default:null" json:"descriptionEn"`
	ResponsibilityRu string     `gorm:"type:text;default:null" json:"responsibilityRu"`
	ResponsibilityUz string     `gorm:"type:text;default:null" json:"responsibilityUz"`
	ResponsibilityEn string     `gorm:"type:text;default:null" json:"responsibilityEn"`
	RequirementRu    string     `gorm:"type:text;default:null" json:"requirementRu"`
	RequirementUz    string     `gorm:"type:text;default:null" json:"requirementUz"`
	RequirementEn    string     `gorm:"type:text;default:null" json:"requirementEn"`
	Region           string     `gorm:"type:varchar(255);default:null" json:"region"`
	TypeUz           string     `gorm:"type:varchar(500)" json:"typeUz"`
	TypeEn           string     `gorm:"type:varchar(500)" json:"typeEn"`
	TypeRu           string     `gorm:"type:varchar(500)" json:"typeRu"`
	Image            string     `gorm:"type:varchar(300)" json:"image"`
	IsActive         *bool      `gorm:"type:boolean;default:true;index" json:"isActive"`
	IsDeleted        *bool      `gorm:"type:boolean;default:false;index" json:"isDeleted"`
	Created          *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID        *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt        *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
	Updated          *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID        *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt        *time.Time `gorm:"type:timestamptz;default:null" json:"updatedAt"`
}

type Applicant struct {
	ID          int        `gorm:"type:bigint;primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(250)" json:"name"`
	Phone       string     `gorm:"type:varchar(250)" json:"phone"`
	Status      int8       `gorm:"type:smallint;default:0" json:"status"`
	Vacancy     *Vacancy   `gorm:"foreignKey:VacancyID;constraint:OnDelete:CASCADE;" json:"vacancy"`
	VacancyID   *int       `gorm:"type:bigint;default:null" json:"vacancyId"`
	Description string     `gorm:"type:varchar(500)" json:"description"`
	Resume      string     `gorm:"type:varchar(300)" json:"resume"`
	CreatedAt   *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
}
type ApplicantFilter struct {
	Status    *int `form:"status" json:"status"`
	Page      int  `form:"page" json:"page"`
	PageSize  int  `form:"pageSize" json:"pageSize"`
	VacancyId int  `form:"vacancyId" json:"vacancyId"`
}

type ApplicantRequest struct {
	Name        string `json:"name" form:"name"`
	Phone       string `json:"phone" form:"phone"`
	VacancyId   int    `json:"vacancyId" form:"vacancyId"`
	Description string `json:"description" form:"description"`
	Resume      string `json:"-" form:"-"`
}

type VacancyRequest struct {
	NameRu           string `json:"nameRu" form:"nameRu"`
	NameEn           string `json:"nameEn" form:"nameEn"`
	NameUz           string `json:"nameUz" form:"nameUz"`
	DescriptionEn    string `json:"descriptionEn" form:"descriptionEn"`
	DescriptionRu    string `json:"descriptionRu" form:"descriptionRu"`
	DescriptionUz    string `json:"descriptionUz" form:"descriptionUz"`
	ResponsibilityEn string `json:"responsibilityEn" form:"responsibilityEn"`
	ResponsibilityRu string `json:"responsibilityRu" form:"responsibilityRu"`
	ResponsibilityUz string `json:"responsibilityUz" form:"responsibilityUz"`
	RequirementEn    string `json:"requirementEn" form:"requirementEn"`
	RequirementRu    string `json:"requirementRu" form:"requirementRu"`
	RequirementUz    string `json:"requirementUz" form:"requirementUz"`
	Region           string `json:"region" form:"region"`
	TypeUz           string `json:"typeUz" form:"typeUz"`
	TypeEn           string `json:"typeEn" form:"typeEn"`
	TypeRu           string `json:"typeRu" form:"typeRu"`
	IsActive         *bool  `json:"isActive" form:"isActive"`
	Image            string `json:"-" form:"-"`
}
