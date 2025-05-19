package subscriptionService

import (
	"context"
	"fmt"
	"sync"

	model "github.com/5aradise/gather-weather/internal/models"
	"github.com/5aradise/gather-weather/internal/models/frequency"
	. "github.com/5aradise/gather-weather/pkg/types"
	"github.com/google/uuid"
)

type (
	service struct {
		stor      subscriptionStorage
		validator validator

		tokens *SyncMap[uuid.UUID, model.Subscription]

		subs struct {
			mu     sync.Mutex
			hourly map[uuid.UUID]model.SubShort
			daily  map[uuid.UUID]model.SubShort
		}
	}

	subscriptionStorage interface {
		ListAllSubscriptions(ctx context.Context) ([]model.Subscription, error)
		CheckSubscriptionByEmail(ctx context.Context, email string) (bool, error)
		CreateSubscription(ctx context.Context, sub model.Subscription) (model.Subscription, error)
		DeleteSubscriptionByToken(ctx context.Context, token uuid.UUID) error
	}

	validator interface {
		ValidateSubscription(model.Subscription) error
	}
)

func New(stor subscriptionStorage, v validator) (*service, error) {
	subs, err := stor.ListAllSubscriptions(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't list subscriptions: %w", err)
	}

	hourly := make(map[uuid.UUID]model.SubShort)
	daily := make(map[uuid.UUID]model.SubShort)
	for _, sub := range subs {
		switch sub.Frequency {
		case frequency.Hourly:
			hourly[sub.Token] = model.SubShort{
				Email: sub.Email,
				City:  sub.City,
			}
		case frequency.Daily:
			daily[sub.Token] = model.SubShort{
				Email: sub.Email,
				City:  sub.City,
			}
		}
	}

	return &service{
		stor:      stor,
		validator: v,
		tokens:    NewSyncMap[uuid.UUID, model.Subscription](),
		subs: struct {
			mu     sync.Mutex
			hourly map[uuid.UUID]model.SubShort
			daily  map[uuid.UUID]model.SubShort
		}{
			hourly: hourly,
			daily:  daily,
		},
	}, nil
}
