package weatherService

import (
	"fmt"
	"strings"

	"github.com/5aradise/gather-weather/config"
	model "github.com/5aradise/gather-weather/internal/models"
	"github.com/valyala/fasthttp"
)

func (s *service) CurrentWeather(city string) (model.Weather, config.ServiceError) {
	city = strings.ToLower(city)

	weather, ok := s.currWeather.Get(city)
	if ok {
		if weather == (model.Weather{}) {
			return model.Weather{}, config.NewServiceErr(
				config.CodeNotFound,
				config.ErrCityNotFound,
			)
		}

		return weather, config.ServiceError{}
	}

	weather, err := s.getCurrentWeather(city)
	s.currWeather.Set(string([]byte(city)), weather)
	if err != nil {
		if err == ErrNoLocMatchq {
			return model.Weather{}, config.NewServiceErr(
				config.CodeNotFound,
				config.ErrCityNotFound,
			)
		}

		return model.Weather{}, config.NewUnknownErr(err)
	}

	return weather, config.ServiceError{}
}

func (s *service) getCurrentWeather(city string) (model.Weather, error) {
	url := s.currWeatherURL(city)
	status, data, err := fasthttp.Get(nil, url)
	if err != nil {
		return model.Weather{}, err
	}
	if status != fasthttp.StatusOK {
		var res errResponse
		err = s.unmarshal(data, &res)
		if err != nil {
			return model.Weather{}, err
		}

		return model.Weather{}, res.Err()
	}

	var res response
	err = s.unmarshal(data, &res)
	if err != nil {
		return model.Weather{}, err
	}

	return model.Weather{
		Temperature: res.Current.TempC,
		Humidity:    res.Current.Humidity,
		Description: res.Current.Condition.Text,
	}, nil
}

func (s *service) currWeatherURL(city string) string {
	const currEndPointTemplate = "%s/current.json?key=%s&q=%s&aqi=no"

	return fmt.Sprintf(currEndPointTemplate, ApiURL, s.key, city)
}
