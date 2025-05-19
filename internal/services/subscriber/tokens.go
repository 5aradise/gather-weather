package subscriptionService

import (
	"time"

	model "github.com/5aradise/gather-weather/internal/models"
	"github.com/google/uuid"
)

func (s *service) generateToken(sub model.Subscription, deadline time.Duration) uuid.UUID {
	token := uuid.New()
	sub.Token = token
	s.tokens.Set(token, sub)
	time.AfterFunc(deadline, func() {
		s.tokens.Delete(token)
	})
	return token
}

func (s *service) pullSubByToken(token uuid.UUID) (model.Subscription, bool) {
	return s.tokens.Pull(token)
}
