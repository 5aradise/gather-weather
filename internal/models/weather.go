package model

import "fmt"

type Weather struct {
	Temperature float32
	Humidity    float32
	Description string
}

func (w Weather) Format() string {
	return fmt.Sprintf(
		"Weather:\r\n"+
			"	Temperature: %.2f🌡\r\n"+
			"	Humidity: %.2f💧\r\n"+
			"	Description: %s\r\n",
		w.Temperature, w.Humidity, w.Description,
	)
}
