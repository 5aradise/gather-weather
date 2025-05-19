package req

import (
	model "github.com/5aradise/gather-weather/internal/models"
	"github.com/5aradise/gather-weather/internal/models/frequency"
	"github.com/gofiber/fiber/v3"
)

func SubscriptionFromForm(c fiber.Ctx) (model.Subscription, error) {
	freq, err := frequency.New(c.FormValue("frequency"))
	if err != nil {
		return model.Subscription{}, err
	}

	return model.Subscription{
		Email:     c.FormValue("email"),
		City:      c.FormValue("city"),
		Frequency: freq,
	}, nil
}
