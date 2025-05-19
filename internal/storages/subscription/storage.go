package subscriptionStorage

import (
	"context"
	"errors"

	"github.com/5aradise/gather-weather/config"
	model "github.com/5aradise/gather-weather/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s storage) CreateSubscription(ctx context.Context, body model.Subscription) (model.Subscription, error) {
	err := s.db.WithContext(ctx).Create(&body).Error
	return body, err
}

func (s storage) ListAllSubscriptions(ctx context.Context) ([]model.Subscription, error) {
	var subs []model.Subscription
	err := s.db.WithContext(ctx).
		Find(&subs).Error
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s storage) CheckSubscriptionByEmail(ctx context.Context, email string) (bool, error) {
	var token model.Subscription
	err := s.db.WithContext(ctx).
		Select("token").
		Where("email", email).
		First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s storage) DeleteSubscriptionByToken(ctx context.Context, token uuid.UUID) error {
	res := s.db.WithContext(ctx).
		Where("token", token).
		Delete(model.Subscription{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return config.ErrTokenNotFound
	}

	return nil
}
