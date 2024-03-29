package models

import "time"

type Customer struct {
	ID        int        `gorm:"type:bigint;primaryKey" json:"id"`
	Name      string     `gorm:"type:varchar(255)" json:"name"`
	Phone     string     `gorm:"type:varchar(255);unique" json:"phone"`
	Password  string     `gorm:"type:varchar(255)" json:"-"`
	Email     string     `gorm:"type:varchar(255)" json:"email"`
	Birthday  string     `gorm:"type:date;default:null" json:"birthday"`
	Created   *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID *int       `gorm:"type:bigint;default:null"  json:"-"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
	Updated   *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID *int       `gorm:"type:bigint;default:null"  json:"-"`
	UpdatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"updated_at"`
	LastVisit *time.Time `gorm:"type:timestamptz;default:null" json:"last_visit"`
}

type CustomerFavorites struct {
	Customer   *Customer `gorm:"foreignKey:CustomerID" json:"customer"`
	CustomerID *int      `gorm:"type:bigint;default:null" json:"-"`
	Product    *Products `gorm:"foreignKey:ProductID" json:"product"`
	ProductID  *int      `gorm:"type:bigint;default:null" json:"-"`
}

type Codes struct {
	ID        uint32     `gorm:"type:bigint;primaryKey" json:"id"`
	Phone     int        `gorm:"type:bigint;index" json:"phone"`
	Code      int        `gorm:"type:integer" json:"code"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
}
type CustomerMetadata struct {
	Id int
}
type CustomerLogin struct {
	Phone    string `json:"phone" example:"998995117361"`
	Password string `json:"password" example:"secret"`
}
type CustomerRegister struct {
	Name     string `json:"name" `
	Password string `json:"password" example:"secret"`
}

type CheckCodeRequest struct {
	Phone        string `json:"phone" example:"998995117361"`
	Code         string `json:"code" example:"997361"`
	IsRegistered bool   `json:"isRegistered"`
}
type CustomerRegisterByAdmin struct {
	Name  string `json:"name" form:"name"`
	Phone string `json:"phone" form:"phone"`
}

type CustomerRequest struct {
	Name     string `json:"name" form:"name" `
	Email    string `json:"email" form:"email"`
	Birthday string `json:"birthday" form:"birthday"`
}
type CustomerRegisterRequest struct {
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password" `
}
type CustomerCode struct {
	Phone string `json:"phone" example:"998995117361"`
}
