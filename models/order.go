package models

import "time"

type Orders struct {
	ID          int        `gorm:"type:bigint not null;primaryKey" json:"id"`
	Customer    *Customer  `gorm:"foreignKey:CustomerID" json:"customer"`
	CustomerID  int        `gorm:"type:bigint;default:null;index" json:"customer_id"`
	Description string     `gorm:"type:varchar(500) not null" json:"description"`
	Status      *int8      `gorm:"type:smallint not null;default:0;index" json:"status"`
	Total       float64    `gorm:"type:decimal(16,2) not null" json:"total"`
	Created     *Admins    `gorm:"foreignKey:CreatedID"       json:"created"`
	CreatedID   *int       `gorm:"type:integer;default:null"  json:"-"`
	CreatedAt   *time.Time `gorm:"type:timestamptz;default:null;index" json:"created_at"`
	Updated     *Admins    `gorm:"foreignKey:UpdatedID"       json:"updated"`
	UpdatedID   *int       `gorm:"type:integer;default:null"  json:"-"`
	UpdatedAt   *time.Time `gorm:"type:timestamptz;default:null" json:"updated_at"`
	DeletedID   *int       `gorm:"type:integer;default:null"  json:"-"`
	Deleted     *Admins    `gorm:"foreignKey:DeletedID"       json:"deleted"`
	DeletedAt   *time.Time `gorm:"type:timestamptz;default:null" json:"deleted_at"`
}

type OrderApplicant struct {
	ID        int        `gorm:"type:bigint not null;primaryKey" json:"id"`
	FullName  string     `gorm:"type:varchar(400)" json:"full_name"`
	Phone     string     `gorm:"type:varchar(20)" json:"phone"`
	Message   string     `gorm:"type:text" json:"message"`
	CreatedAt *time.Time `gorm:"type:timestamptz;default:null" json:"created_at"`
}

type OrderApplicantRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Message  string `json:"message"`
}
type OrderItems struct {
	Price  float64   `gorm:"type:decimal(16,2) not null" json:"total"`
	Amount int       `gorm:"type:integer;default:null" json:"amount"`
	Item   *Products `gorm:"foreignKey:ItemId" json:"item"`
	ItemId int       `gorm:"type:integer;default:null" json:"item_id"`
}

type OrderResponse struct {
	*Orders
	Items []OrderItems `json:"items"`
}

type OrderApplicantResponse struct {
	Page      int              `json:"page"`
	PageSize  int              `json:"page_size"`
	Count     int              `json:"count"`
	Applicant []OrderApplicant `json:"orderApplicants"`
}

type CustomOrderResponse struct {
	Orders         []Orders `json:"orders"`
	Page           int      `json:"page"`
	PageSize       int      `json:"page_size"`
	Count          int      `json:"count"`
	ActiveCount    int      `json:"active_count"`
	FinishedCount  int      `json:"finished_count"`
	CancelledCount int      `json:"cancelled_count"`
}

type OrderRequest struct {
	Description string              `json:"description"`
	Total       float64             `json:"total"`
	Items       []OrderItemsRequest `json:"items"`
}
type OrderUpdateRequest struct {
	CustomerId  int                 `json:"customer_id" form:"customer_id"`
	Description string              `json:"description" form:"description"`
	Total       *float64            `json:"total" form:"total"`
	Items       []OrderItemsRequest `json:"items" form:"items"`
}

type AdminOrderFilter struct {
	CustomerID int    `json:"customer_id"`
	DateFrom   string `json:"date_from" form:"date_from"`
	Status     *int   `json:"status" form:"status"`
	DateTo     string `json:"date_to" form:"date_to"`
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"page_size" form:"page_size"`
}

type OrderApplicantFilter struct {
	Phone    string `json:"phone" form:"phone"`
	FullName string `json:"full_name" form:"full_name"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
}

type OrderFilter struct {
	Status   *int `json:"status" form:"status"`
	Page     int  `json:"page" form:"page"`
	PageSize int  `json:"page_size" form:"page_size"`
}

type OrderItemsRequest struct {
	Price  float64 `json:"price" form:"price"`
	Amount int     `json:"amount" form:"amount"`
	ItemID int     `json:"item_id" form:"item_id"`
}
