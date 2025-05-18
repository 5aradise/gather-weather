package config

import (
	"log"
)

type ServiceError struct {
	ServiceCode ServiceCode
	Err         error
}

func (e ServiceError) Error() string {
	if e.Err == nil {
		return ""
	}

	return e.Err.Error()
}

func (e ServiceError) Unwrap() error {
	return e.Err
}

func (e ServiceError) IsZero() bool {
	return e.ServiceCode == 0 && e.Err == nil
}

func NewServiceErr(code ServiceCode, err error) ServiceError {
	serr := ServiceError{
		ServiceCode: code,
		Err:         err,
	}
	return serr
}

func NewUnknownErr(err error) ServiceError {
	log.Println("UNKNOWN ERROR: ", err)

	return ServiceError{
		Err: err,
	}
}
