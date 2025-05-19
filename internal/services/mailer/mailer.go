package mailService

import (
	"fmt"
	"net/smtp"

	"github.com/5aradise/gather-weather/config"
)

const template = "To: %s\r\n" +
	"Subject: %s\r\n" +
	"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
	"\r\n" +
	"%s\r\n"

func (s *service) SendMail(to, subject, message string) config.ServiceError {
	err := smtp.SendMail(s.addr, s.auth, s.sender, []string{to}, fmt.Appendf(nil,
		template, to, subject, message,
	))
	if err != nil {
		return config.NewUnknownErr(err)
	}

	return config.ServiceError{}
}
