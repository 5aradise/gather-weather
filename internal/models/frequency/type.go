package frequency

import (
	"errors"
	"time"
)

type Type string

const (
	Hourly Type = "hourly"
	Daily  Type = "daily"
)

func New(s string) (Type, error) {
	switch s {
	case "hourly":
		return Hourly, nil
	case "daily":
		return Daily, nil
	default:
		return Type(""), errors.New("frequency: unknown duration")
	}
}

func (t Type) Duration() time.Duration {
	switch t {
	case Hourly:
		return time.Hour
	case Daily:
		return 24 * time.Hour
	default:
		panic("unknown duration")
	}
}
