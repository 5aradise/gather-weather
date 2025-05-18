package res

import model "github.com/5aradise/gather-weather/internal/models"

type Weather struct {
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	Description string  `json:"description"`
}

func ModelToWeather(m model.Weather) Weather {
	return Weather{
		Temperature: m.Temperature,
		Humidity:    m.Humidity,
		Description: m.Description,
	}
}
