package models

import "time"

type Country struct {
	ID        int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz    string     `gorm:"type:varchar(250) not null" json:"nameUz"`
	NameRu    string     `gorm:"type:varchar(250) not null" json:"nameRu"`
	NameEn    string     `gorm:"type:varchar(250) not null" json:"nameEn"`
	Created   *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
}

type CountryRequest struct {
	NameUz string `json:"nameUz" form:"nameUz"`
	NameRu string `json:"nameRu" form:"nameRu"`
	NameEn string `json:"nameEn" form:"nameEn"`
}
