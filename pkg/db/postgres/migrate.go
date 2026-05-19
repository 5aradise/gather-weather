package postgres

import (
	model "github.com/5aradise/gather-weather/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Subscription{})
}
