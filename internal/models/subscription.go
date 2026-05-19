package model

import (
	"github.com/5aradise/gather-weather/internal/models/frequency"
	"github.com/google/uuid"
)

type Subscription struct {
	Token     uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Email     string         `gorm:"type:text;not null;uniqueIndex" validate:"required,email"`
	City      string         `gorm:"type:text;not null" validate:"required"`
	Frequency frequency.Type `gorm:"type:frequency_type;not null" validate:"required"`
}

type SubShort struct {
	Email string
	City  string
}

func (Subscription) TableName() string {
	return "subscriptions"
}
