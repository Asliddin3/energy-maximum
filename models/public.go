package models

import "time"

type PublicOffer struct {
	ID            int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz        string     `gorm:"type:varchar(250);default:null" json:"name_uz"`
	NameRu        string     `gorm:"type:varchar(250);default:null" json:"name_ru"`
	NameEn        string     `gorm:"type:varchar(250);default:null" json:"name_en"`
	DescriptionUz string     `gorm:"type:text;default:null" json:"description_uz"`
	DescriptionRu string     `gorm:"type:text;default:null" json:"description_ru"`
	DescriptionEn string     `gorm:"type:text;default:null" json:"description_en"`
	Created       *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID     *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt     *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
}

type PublicOfferRequest struct {
	NameUz        string `json:"name_uz" form:"name_uz"`
	NameRu        string `json:"name_ru" form:"name_ru"`
	NameEn        string `json:"name_en" form:"name_en"`
	DescriptionEn string `json:"description_en" form:"description_en"`
	DescriptionRu string `json:"description_ru" form:"description_ru"`
	DescriptionUz string `json:"description_uz" form:"description_uz"`
}
