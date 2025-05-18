package weatherHandler

import (
	"github.com/5aradise/gather-weather/config"
	model "github.com/5aradise/gather-weather/internal/models"
	"github.com/gofiber/fiber/v3"
)

type (
	handler struct {
		srv weatherer
	}

	weatherer interface {
		CurrentWeather(city string) (model.Weather, config.ServiceError)
	}
)

func New(w weatherer) *handler {
	return &handler{w}
}

func (h *handler) Init(r fiber.Router) {
	r.Get("/weather", h.getCurrentWeather)
}
