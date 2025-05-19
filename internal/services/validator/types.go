package validationServ

import "github.com/go-playground/validator/v10"

type (
	service struct {
		pgv *validator.Validate

		checkCity func(city string) bool
	}
)

func New(checkCity func(city string) bool) service {
	return service{
		pgv:       validator.New(),
		checkCity: checkCity,
	}
}
