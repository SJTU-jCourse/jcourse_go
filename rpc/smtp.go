package rpc

import (
	"context"

	"gopkg.in/gomail.v2"

	"jcourse_go/config"
)

type EmailSender interface {
	SendMail(ctx context.Context, recipient string, subject string, body string) error
}

type NilSender struct {
}

func (s *NilSender) SendMail(ctx context.Context, recipient string, subject string, body string) error {
	return nil
}

type SMTPSender struct {
	config config.SMTP
}

func (s *SMTPSender) SendMail(ctx context.Context, recipient string, subject string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.config.Sender)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", subject)
	// text/html for a html email
	msg.SetBody("text/plain", body)

	n := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)

	// Send the email
	err := n.DialAndSend(msg)
	return err
}
