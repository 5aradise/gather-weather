package weatherService

func (s *service) CheckCity(city string) bool {
	_, serr := s.CurrentWeather(city)
	return serr.IsZero()
}
