package models

import "time"

type Vacancy struct {
	ID               int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz           string     `gorm:"type:varchar(250) not null" json:"name_uz"`
	NameRu           string     `gorm:"type:varchar(250) not null" json:"name_ru"`
	NameEn           string     `gorm:"type:varchar(250) not null" json:"name_en"`
	DescriptionRu    string     `gorm:"type:text;default:null" json:"description_ru"`
	DescriptionUz    string     `gorm:"type:text;default:null" json:"description_uz"`
	DescriptionEn    string     `gorm:"type:text;default:null" json:"description_en"`
	ResponsibilityRu string     `gorm:"type:text;default:null" json:"responsibility_ru"`
	ResponsibilityUz string     `gorm:"type:text;default:null" json:"responsibility_uz"`
	ResponsibilityEn string     `gorm:"type:text;default:null" json:"responsibility_en"`
	RequirementRu    string     `gorm:"type:text;default:null" json:"requirement_ru"`
	RequirementUz    string     `gorm:"type:text;default:null" json:"requirement_uz"`
	RequirementEn    string     `gorm:"type:text;default:null" json:"requirement_en"`
	Region           string     `gorm:"type:varchar(255);default:null" json:"region"`
	TypeUz           string     `gorm:"type:varchar(500)" json:"type_uz"`
	TypeEn           string     `gorm:"type:varchar(500)" json:"type_en"`
	TypeRu           string     `gorm:"type:varchar(500)" json:"type_ru"`
	Image            string     `gorm:"type:varchar(300)" json:"image"`
	IsActive         *bool      `gorm:"type:boolean;default:true;index" json:"is_active"`
	IsDeleted        *bool      `gorm:"type:boolean;default:false;index" json:"is_deleted"`
	Created          *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID        *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt        *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
	Updated          *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID        *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt        *time.Time `gorm:"type:timestamptz;default:null" json:"updated_at"`
}

type Applicant struct {
	ID          int        `gorm:"type:bigint;primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(250)" json:"name"`
	Phone       string     `gorm:"type:varchar(250)" json:"phone"`
	Status      int8       `gorm:"type:smallint;default:0" json:"status"`
	Vacancy     *Vacancy   `gorm:"foreignKey:VacancyID;constraint:OnDelete:CASCADE;" json:"vacancy"`
	VacancyID   *int       `gorm:"type:bigint;default:null" json:"vacancy_id"`
	Description string     `gorm:"type:varchar(500)" json:"description"`
	Resume      string     `gorm:"type:varchar(300)" json:"resume"`
	CreatedAt   *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
}
type ApplicantFilter struct {
	Status    *int `form:"status" json:"status"`
	Page      int  `form:"page" json:"page"`
	PageSize  int  `form:"page_size" json:"page_size"`
	VacancyId int  `form:"vacancy_id" json:"vacancy_id"`
}

type ApplicantRequest struct {
	Name        string `json:"name" form:"name"`
	Phone       string `json:"phone" form:"phone"`
	VacancyId   int    `json:"vacancy_id" form:"vacancy_id"`
	Description string `json:"description" form:"description"`
	Resume      string `json:"-" form:"-"`
}

type VacancyRequest struct {
	NameRu           string `json:"name_ru" form:"name_ru"`
	NameEn           string `json:"name_en" form:"name_en"`
	NameUz           string `json:"name_uz" form:"name_uz"`
	DescriptionEn    string `json:"description_en" form:"description_en"`
	DescriptionRu    string `json:"description_ru" form:"description_ru"`
	DescriptionUz    string `json:"description_uz" form:"description_uz"`
	ResponsibilityEn string `json:"responsibility_en" form:"responsibility_en"`
	ResponsibilityRu string `json:"responsibility_ru" form:"responsibility_ru"`
	ResponsibilityUz string `json:"responsibility_uz" form:"responsibility_uz"`
	RequirementEn    string `json:"requirement_en" form:"requirement_en"`
	RequirementRu    string `json:"requirement_ru" form:"requirement_ru"`
	RequirementUz    string `json:"requirement_uz" form:"requirement_uz"`
	Region           string `json:"region" form:"region"`
	TypeUz           string `json:"type_uz" form:"type_uz"`
	TypeEn           string `json:"type_en" form:"type_en"`
	TypeRu           string `json:"type_ru" form:"type_ru"`
	IsActive         *bool  `json:"is_active" form:"is_active"`
	Image            string `json:"-" form:"-"`
}
