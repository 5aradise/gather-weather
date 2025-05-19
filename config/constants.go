package config

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

var ( // Errors
	ErrCityNotFound    = errors.New("city not found")
	ErrEmailSubscribed = errors.New("email already subscribed")
	ErrTokenNotFound   = errors.New("token not found")
	ErrTODO            = errors.New("TODO")
)

type ServiceCode int

const (
	CodeBadRequest ServiceCode = iota + 1
	CodeEmptyValue
	CodeConflict
	CodeUnprocessableEntity

	CodeNotFound

	CodeUnauthorized
	CodeForbidden
)

func (code ServiceCode) ToHttpStatus() int {
	switch code {
	case CodeBadRequest, CodeEmptyValue:
		return fiber.StatusBadRequest
	case CodeUnauthorized:
		return fiber.StatusUnauthorized
	case CodeForbidden:
		return fiber.StatusForbidden
	case CodeNotFound:
		return fiber.StatusNotFound
	case CodeConflict:
		return fiber.StatusConflict
	case CodeUnprocessableEntity:
		return fiber.StatusUnprocessableEntity
	default:
		return fiber.StatusInternalServerError
	}
}
