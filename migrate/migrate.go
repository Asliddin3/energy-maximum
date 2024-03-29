package migrate

import (
	"github.com/Asliddin3/energy-maximum/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Admins{},
		&models.PublicOffer{},
		&models.Products{},
		&models.Brand{},
		&models.Customer{},
		&models.Pages{},
		&models.Roles{},
		&models.RoleItems{},
		&models.ModuleItems{},
		&models.Modules{},
		&models.CustomerFavorites{},
		&models.Country{},
		&models.Contact{},
		&models.Codes{},
	)
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&models.About{},
		&models.Orders{},
		&models.OrderItems{},
		&models.Service{},
		&models.ProductAdditions{},
		&models.Analog{},
		&models.AnalogProduct{},
		&models.Banner{},
		&models.Category{},
		&models.ProductMedia{},
		&models.Applicant{},
		&models.OrderApplicant{},
	)
	if err != nil {
		return err
	}
	err = db.AutoMigrate(
		&models.Vacancy{},
		&models.News{},
		&models.Parameters{},
		&models.ProductParameters{},
	)
	if err != nil {
		return err
	}
	err = db.Exec("alter table product_additions add unique  (product_category_id,addition_category_id)").Error
	if err != nil {
		return err
	}
	return nil
}
