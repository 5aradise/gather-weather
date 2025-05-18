package weatherService

import (
	"errors"
	"fmt"
)

var (
	ErrNoLocMatchq = errors.New("no location found matching parameter 'q'")
	ErrBadApiKey   = errors.New("bad api key")
)

type (
	response struct {
		Current struct {
			TempC     float32 `json:"temp_c"`
			Humidity  float32 `json:"humidity"`
			Condition struct {
				Text string `json:"text"`
			} `json:"condition"`
		} `json:"current"`
	}

	errResponse struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
)

func (res errResponse) Err() error {
	switch res.Error.Code {
	case 1006:
		return ErrNoLocMatchq
	default:
		if code := res.Error.Code; 2000 <= code && code < 3000 {
			return fmt.Errorf("%w: %s", ErrBadApiKey, res.Error.Message)
		}
	}
	return errors.New(res.Error.Message)
}
