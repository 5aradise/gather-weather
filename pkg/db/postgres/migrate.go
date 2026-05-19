package postgres

import (
	model "github.com/5aradise/gather-weather/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.Exec(`
DO $$
BEGIN
	CREATE TYPE frequency_type AS ENUM ('hourly', 'daily');
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;
`).Error; err != nil {
		return err
	}

	return db.AutoMigrate(&model.Subscription{})
}
