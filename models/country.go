package models

import "time"

type Country struct {
	ID        int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz    string     `gorm:"type:varchar(250) not null" json:"name_uz"`
	NameRu    string     `gorm:"type:varchar(250) not null" json:"name_ru"`
	NameEn    string     `gorm:"type:varchar(250) not null" json:"name_en"`
	Created   *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
}

type CountryRequest struct {
	NameUz string `json:"name_uz" form:"name_uz"`
	NameRu string `json:"name_ru" form:"name_ru"`
	NameEn string `json:"name_en" form:"name_en"`
}
