package subscriptionStorage

import "gorm.io/gorm"

type (
	storage struct {
		db *gorm.DB
	}
)

func New(db *gorm.DB) storage {
	return storage{db}
}
