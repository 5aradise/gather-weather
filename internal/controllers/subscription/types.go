package subscriptionHandler

import (
	"context"
	"iter"

	"github.com/5aradise/gather-weather/config"
	model "github.com/5aradise/gather-weather/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type (
	handler struct {
		sub     subscriber
		mail    mailer
		weather weatherer
	}

	subscriber interface {
		RequestSubscription(ctx context.Context, sub model.Subscription) (token uuid.UUID, err config.ServiceError)
		ConfirmSubscription(ctx context.Context, token uuid.UUID) config.ServiceError
		Unsubscribe(ctx context.Context, token uuid.UUID) config.ServiceError

		ListHourlySubscribers() iter.Seq[model.SubShort]
		ListDailySubscribers() iter.Seq[model.SubShort]
	}

	weatherer interface {
		CurrentWeather(city string) (model.Weather, config.ServiceError)
	}

	mailer interface {
		SendMail(to, subject, message string) config.ServiceError
	}
)

func New(s subscriber, mail mailer, weather weatherer) *handler {
	return &handler{s, mail, weather}
}

func (h *handler) Init(r fiber.Router) {
	r.Post("/subscribe", h.subscribe)
	r.Get("/confirm/:token", h.confirm)
	r.Get("/unsubscribe/:token", h.unsubscribe)
}
