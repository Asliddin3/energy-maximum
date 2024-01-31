package models

import "time"

type Analog struct {
	ID        int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz    string     `gorm:"type:varchar(250);default:null" json:"nameUz"`
	NameRu    string     `gorm:"type:varchar(250);default:null" json:"nameRu"`
	NameEn    string     `gorm:"type:varchar(250);default:null" json:"nameEn"`
	Created   *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"createdAt"`
	Updated   *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"updatedAt"`
}

type AnalogProduct struct {
	Analog    *Analog   `gorm:"foreignKey:AnalogID;constraint:OnDelete:CASCADE;" json:"analog"`
	AnalogID  int       `gorm:"type:integer not null" json:"analogId"`
	Product   *Products `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"product"`
	ProductID int       `gorm:"type:integer not null" json:"productId"`
}


type AnalogResponse struct {
	*Analog
	Items []Products `json:"items"`
}

type AnalogRequest struct {
	NameUz string `json:"nameUz" form:"nameUz"`
	NameRu string `json:"nameRu" form:"nameRu"`
	NameEn string `json:"nameEn" form:"nameEn"`
}
