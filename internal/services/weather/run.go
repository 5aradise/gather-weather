package weatherService

import (
	"time"
)

const weatherUpdateDelayInMin = 15

func (s *service) run() {
	now := time.Now().Add(time.Second)
	minOff := now.Minute()
	delayToNext := weatherUpdateDelayInMin - minOff%weatherUpdateDelayInMin
	nextReset := now.Add(time.Duration(delayToNext) * time.Minute)

	y, mon, d := nextReset.Date()
	nextReset = time.Date(y, mon, d, nextReset.Hour(), nextReset.Minute(), 0, 0, nextReset.Location())
	time.AfterFunc(nextReset.Sub(time.Now()), func() {
		go s.currWeather.Reset()
		for range time.NewTicker(weatherUpdateDelayInMin * time.Minute).C {
			s.currWeather.Reset()
		}
	})
}
