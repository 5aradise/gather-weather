package res

import (
	"github.com/5aradise/gather-weather/config"
	"github.com/gofiber/fiber/v2"
)

func ServiceErr(c *fiber.Ctx, serr config.ServiceError) error {
	return c.Status(serr.ServiceCode.ToHttpStatus()).SendString(serr.Error())
}
