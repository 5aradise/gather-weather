package subscriptionService

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/5aradise/gather-weather/config"
	model "github.com/5aradise/gather-weather/internal/models"
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

	return config.ServiceError{}
}
