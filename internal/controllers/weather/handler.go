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

	weather, serr := h.srv.CurrentWeather(city)
	if !serr.IsZero() {
		return res.ServiceErr(c, serr)
	}

	return c.Status(fasthttp.StatusOK).
		JSON(res.ModelToWeather(weather))
}
