package model

import (
	"github.com/5aradise/gather-weather/internal/models/frequency"
	"github.com/google/uuid"
)

type Subscription struct {
	Token     uuid.UUID
	Email     string         `validate:"required,email"`
	City      string         `validate:"required"`
	Frequency frequency.Type `validate:"required"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
