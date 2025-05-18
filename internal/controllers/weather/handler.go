package weatherHandler

import (
	res "github.com/5aradise/gather-weather/internal/controllers/response"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
)

func (h *handler) getCurrentWeather(c fiber.Ctx) error {
	city := c.Query("city")
	if city == "" {
		return c.SendStatus(fasthttp.StatusBadRequest)
	}

	weather, err := h.srv.CurrentWeather(city)
	if !err.IsZero() {
		return c.SendStatus(err.ServiceCode.ToHttpStatus())
	}

	return c.Status(fasthttp.StatusOK).
		JSON(res.ModelToWeather(weather))
}
