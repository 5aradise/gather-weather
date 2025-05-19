package subscriptionHandler

import (
	"context"

	"github.com/5aradise/gather-weather/config"
	model "github.com/5aradise/gather-weather/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type (
	handler struct {
		srv subscriber
	}

	subscriber interface {
		RequestSubscription(ctx context.Context, sub model.Subscription) (token uuid.UUID, err config.ServiceError)
		ConfirmSubscription(ctx context.Context, token uuid.UUID) config.ServiceError
		Unsubscribe(ctx context.Context, token uuid.UUID) config.ServiceError
	}
)

func New(s subscriber) *handler {
	return &handler{s}
}

func (h *handler) Init(r fiber.Router) {
	r.Post("/subscribe", h.subscribe)
	r.Get("/confirm/:token", h.confirm)
	r.Get("/unsubscribe/:token", h.unsubscribe)
}
