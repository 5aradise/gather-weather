package subscriptionHandler

import (
	"fmt"

	req "github.com/5aradise/gather-weather/internal/controllers/request"
	res "github.com/5aradise/gather-weather/internal/controllers/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func (h *handler) subscribe(c fiber.Ctx) error {
	sub, err := req.SubscriptionFromForm(c)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	token, serr := h.sub.RequestSubscription(c.Context(), sub)
	if !serr.IsZero() {
		return res.ServiceErr(c, serr)
	}

	serr = h.mail.SendMail(sub.Email, "Gather-weather token", fmt.Sprintf("Your token: %s", token))
	if !serr.IsZero() {
		return res.ServiceErr(c, serr)
	}

	return c.SendStatus(fasthttp.StatusOK)
}

func (h *handler) confirm(c fiber.Ctx) error {
	strToken := c.Params("token")
	if strToken == "" {
		panic("emrty token parameter")
	}

	token, err := uuid.Parse(strToken)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	serr := h.sub.ConfirmSubscription(c.Context(), token)
	if !serr.IsZero() {
		return res.ServiceErr(c, serr)
	}

	return c.SendStatus(fasthttp.StatusOK)
}

func (h *handler) unsubscribe(c fiber.Ctx) error {
	strToken := c.Params("token")
	if strToken == "" {
		panic("emrty token parameter")
	}

	token, err := uuid.Parse(strToken)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	serr := h.sub.Unsubscribe(c.Context(), token)
	if !serr.IsZero() {
		return res.ServiceErr(c, serr)
	}

	return c.SendStatus(fasthttp.StatusOK)
}
