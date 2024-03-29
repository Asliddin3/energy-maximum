package models

import "time"

type Analog struct {
	ID        int        `gorm:"type:bigint;primaryKey" json:"id"`
	NameUz    string     `gorm:"type:varchar(250);default:null" json:"name_uz"`
	NameRu    string     `gorm:"type:varchar(250);default:null" json:"name_ru"`
	NameEn    string     `gorm:"type:varchar(250);default:null" json:"name_en"`
	Created   *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
	Updated   *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"updated_at"`
}

type AnalogProduct struct {
	Analog    *Analog   `gorm:"foreignKey:AnalogID;constraint:OnDelete:CASCADE;" json:"analog"`
	AnalogID  int       `gorm:"type:integer not null" json:"analog_id"`
	Product   *Products `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"product"`
	ProductID int       `gorm:"type:integer not null" json:"product_id"`
}

type AnalogResponse struct {
	*Analog
	Items []Products `json:"items"`
}

type AnalogRequest struct {
	NameUz string `json:"name_uz" form:"name_uz"`
	NameRu string `json:"name_ru" form:"name_ru"`
	NameEn string `json:"name_en" form:"name_en"`
}
