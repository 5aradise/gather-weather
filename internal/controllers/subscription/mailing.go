package subscriptionHandler

import (
	"log"
	"time"
)

const day = 24 * time.Hour

func (h *handler) RunMailing() {
	now := time.Now().Add(time.Second)

	nextHourlyMailing := now.Add(time.Hour)
	y, mon, d := nextHourlyMailing.Date()
	nextHourlyMailing = time.Date(y, mon, d, nextHourlyMailing.Hour(), 0, 0, 0, now.Location())
	time.AfterFunc(time.Until(nextHourlyMailing), func() {
		go h.hourlyMailing()
		for range time.NewTicker(time.Hour).C {
			h.hourlyMailing()
		}
	})

	nextDailyMailing := now.Add(day)
	y, mon, d = nextDailyMailing.Date()
	nextDailyMailing = time.Date(y, mon, d, 0, 0, 0, 0, now.Location())
	time.AfterFunc(time.Until(nextDailyMailing), func() {
		go h.dailyMailing()
		for range time.NewTicker(day).C {
			h.dailyMailing()
		}
	})
}

func (h *handler) hourlyMailing() {
	for sub := range h.sub.ListHourlySubscribers() {
		weather, serr := h.weather.CurrentWeather(sub.City)
		if !serr.IsZero() {
			log.Println("UNKNOWN ERROR: ", serr)
			continue
		}
		serr = h.mail.SendMail(sub.Email, "Hourly weather update", weather.Format())
		if !serr.IsZero() {
			log.Println("UNKNOWN ERROR: ", serr)
		}
	}
}

func (h *handler) dailyMailing() {
	for sub := range h.sub.ListDailySubscribers() {
		weather, serr := h.weather.CurrentWeather(sub.City)
		if !serr.IsZero() {
			log.Println("UNKNOWN ERROR: ", serr)
			continue
		}
		serr = h.mail.SendMail(sub.Email, "Daily weather update", weather.Format())
		if !serr.IsZero() {
			log.Println("UNKNOWN ERROR: ", serr)
		}
	}
}
