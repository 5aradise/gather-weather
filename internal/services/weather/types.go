package weatherService

import (
	"encoding/json"

	model "github.com/5aradise/gather-weather/internal/models"
	. "github.com/5aradise/gather-weather/pkg/types"
)

const ApiURL = "http://api.weatherapi.com/v1"

type service struct {
	key         string
	unmarshal   func([]byte, any) error
	currWeather *SyncMap[string, model.Weather]
}

func New(apiKey string, unmarshaler ...func([]byte, any) error) (*service, error) {
	unmarshal := json.Unmarshal
	if len(unmarshaler) > 0 {
		unmarshal = unmarshaler[0]
	}

	srv := &service{
		key:         apiKey,
		unmarshal:   unmarshal,
		currWeather: NewSyncMap[string, model.Weather](),
	}
	_, err := srv.getCurrentWeather("Kyiv")
	if err != nil {
		return nil, err
	}

	srv.run()

	return srv, nil
}
