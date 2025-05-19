package mailService

import (
	"net"
	"net/smtp"
)

type (
	service struct {
		addr   string
		sender string
		auth   smtp.Auth
	}
)

func New(host, port, senderMail, password string) *service {
	auth := smtp.PlainAuth("", senderMail, password, host)
	return &service{
		addr:   net.JoinHostPort(host, port),
		sender: senderMail,
		auth:   auth,
	}
}
