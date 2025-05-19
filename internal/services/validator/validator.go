package validationServ

import (
	"errors"

	model "github.com/5aradise/gather-weather/internal/models"
)

func (s service) ValidateSubscription(sub model.Subscription) error {
	err := s.pgv.Struct(sub)
	if err != nil {
		return err
	}

	if !s.checkCity(sub.City) {
		return errors.New("")
	}
	return nil
}
