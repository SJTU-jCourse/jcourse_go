package mail

import (
	"context"

	"gopkg.in/gomail.v2"

	"jcourse_go/config"
)

type EmailSender interface {
	SendMail(ctx context.Context, recipient string, subject string, body string) error
}

func NewSMTPSender(conf *config.SMTP) EmailSender {
	return &SMTPSender{
		host:     conf.Host,
		port:     conf.Port,
		sender:   conf.Sender,
		username: conf.Username,
		password: conf.Password,
	}
}

type SMTPSender struct {
	host     string
	port     int
	sender   string
	username string
	password string
}

func (s *SMTPSender) SendMail(ctx context.Context, recipient string, subject string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.sender)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", subject)
	// text/html for a html email
	msg.SetBody("text/plain", body)

	n := gomail.NewDialer(s.host, s.port, s.username, s.password)

	// Send the email
	err := n.DialAndSend(msg)
	return err
}
