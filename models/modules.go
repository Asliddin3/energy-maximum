package models

type Modules struct {
	ID          uint32 `gorm:"type:bigint;primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(300);unique;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}

type CreateModuleInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type UpdateModuleInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModuleItems struct {
	ID          uint32   `gorm:"integer;autoIncrement:true" json:"id"`
	Module      *Modules `gorm:"foreignKey:ModuleID;onDelete:CASCADE" json:"module"`
	Name        string   `gorm:"type:varchar(300) not null" json:"name"`
	EndPoint    string   `gorm:"type:varchar(400) not null" json:"end_point"`
	Method      string   `gorm:"type:varchar(100) not null;index" json:"method"`
	Key         string   `gorm:"type:varchar(255);primaryKey" json:"key"`
	Description string   `gorm:"type:text" json:"description"`
	ModuleID    uint32   `gorm:"type:integer ;not null;index" json:"-"`
}

type CreateModuleItemInput struct {
	EndPoint    string `json:"end_point" binding:"required"`
	Method      string `json:"method" binding:"required"`
	Key         string `json:"key" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	ModuleID    uint32 `json:"module_id" binding:"required"`
}

type UpdateModuleItemInput struct {
	EndPoint    string `json:"end_point" binding:"required"`
	Method      string `json:"method" binding:"required"`
	Key         string `json:"key" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	ModuleID    uint32 `json:"module_id" binding:"required"`
}

type GetModuleItemFilter struct {
	ModuleID *uint32 `form:"module_id"`
}
