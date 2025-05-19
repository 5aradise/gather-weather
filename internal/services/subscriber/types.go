package subscriptionService

import (
	"context"

	model "github.com/5aradise/gather-weather/internal/models"
	. "github.com/5aradise/gather-weather/pkg/types"
	"github.com/google/uuid"
)

type (
	service struct {
		stor      subscriptionStorage
		validator validator

		tokens *SyncMap[uuid.UUID, model.Subscription]
	}

	subscriptionStorage interface {
		CheckSubscriptionByEmail(ctx context.Context, email string) (bool, error)
		CreateSubscription(ctx context.Context, sub model.Subscription) (model.Subscription, error)
		DeleteSubscriptionByToken(ctx context.Context, token uuid.UUID) error
	}

	validator interface {
		ValidateSubscription(model.Subscription) error
	}
)

func New(stor subscriptionStorage, v validator) *service {
	return &service{
		stor:      stor,
		validator: v,
		tokens:    NewSyncMap[uuid.UUID, model.Subscription](),
	}
}
