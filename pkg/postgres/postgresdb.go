package postgresdb

import (
	"fmt"

	"github.com/Asliddin3/energy-maximum/config"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewClient(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Tashkent",
		cfg.DB_HOST, cfg.DB_USER, cfg.DB_PASS, cfg.DB_NAME, cfg.DB_PORT)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable:       true,
			IdentifierMaxLength: 0,
		},
	})
	if err := db.Use(otelgorm.NewPlugin()); err != nil {

		panic(err)
	}

	// db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	// db.Set("gorm:table_options", "ENGINE=InnoDB")

	return db, err
}
