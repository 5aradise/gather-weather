package subscriptionService

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"time"

	"github.com/5aradise/gather-weather/config"
	model "github.com/5aradise/gather-weather/internal/models"
	"github.com/5aradise/gather-weather/internal/models/frequency"
	"github.com/google/uuid"
)

const confirmTokenDedline = 5 * time.Minute

func (s *service) RequestSubscription(ctx context.Context, sub model.Subscription) (uuid.UUID, config.ServiceError) {
	err := s.validator.ValidateSubscription(sub)
	if err != nil {
		return uuid.Nil, config.NewServiceErr(
			config.CodeBadRequest,
			fmt.Errorf("invalid subscription: %w", err),
		)
	}

	ok, err := s.stor.CheckSubscriptionByEmail(ctx, sub.Email)
	if err != nil {
		return uuid.Nil, config.NewUnknownErr(err)
	}
	if ok {
		return uuid.Nil, config.NewServiceErr(
			config.CodeConflict,
			config.ErrEmailSubscribed,
		)
	}

	token := s.generateToken(sub, confirmTokenDedline)
	return token, config.ServiceError{}
}

func (s *service) ConfirmSubscription(ctx context.Context, token uuid.UUID) config.ServiceError {
	sub, ok := s.pullSubByToken(token)
	if !ok {
		return config.NewServiceErr(
			config.CodeNotFound,
			config.ErrTokenNotFound,
		)
	}

	_, err := s.stor.CreateSubscription(ctx, sub)
	if err != nil {
		return config.NewServiceErr(
			config.CodeConflict,
			err,
		)
	}

	s.subs.mu.Lock()
	switch sub.Frequency {
	case frequency.Hourly:
		s.subs.hourly[sub.Token] = model.SubShort{
			Email: sub.Email,
			City:  sub.City,
		}
	case frequency.Daily:
		s.subs.daily[sub.Token] = model.SubShort{
			Email: sub.Email,
			City:  sub.City,
		}
	}
	s.subs.mu.Unlock()

	return config.ServiceError{}
}

func (s *service) Unsubscribe(ctx context.Context, token uuid.UUID) config.ServiceError {
	err := s.stor.DeleteSubscriptionByToken(ctx, token)
	if err != nil {
		if errors.Is(err, config.ErrTokenNotFound) {
			return config.NewServiceErr(
				config.CodeNotFound,
				err,
			)
		}

		return config.NewUnknownErr(err)
	}

	s.subs.mu.Lock()
	delete(s.subs.hourly, token)
	delete(s.subs.daily, token)
	s.subs.mu.Unlock()

	return config.ServiceError{}
}

func (s *service) ListHourlySubscribers() iter.Seq[model.SubShort] {
	return func(yield func(model.SubShort) bool) {
		s.subs.mu.Lock()
		defer s.subs.mu.Unlock()
		for _, sub := range s.subs.hourly {
			if !yield(sub) {
				return
			}
		}
	}
}

func (s *service) ListDailySubscribers() iter.Seq[model.SubShort] {
	return func(yield func(model.SubShort) bool) {
		s.subs.mu.Lock()
		defer s.subs.mu.Unlock()
		for _, sub := range s.subs.daily {
			if !yield(sub) {
				return
			}
		}
	}
}
