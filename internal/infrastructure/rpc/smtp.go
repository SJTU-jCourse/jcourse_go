package rpc

import (
	"context"

	"gopkg.in/gomail.v2"

	"jcourse_go/internal/config"
	"jcourse_go/internal/domain/email"
)

type SMTPEmailSender struct {
	conf config.SMTPConfig
}

func NewSMTPEmailSender(conf config.SMTPConfig) email.EmailSender {
	return &SMTPEmailSender{conf: conf}
}

func (s *SMTPEmailSender) SendEmail(ctx context.Context, req email.Request) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.conf.Sender)
	msg.SetHeader("To", req.Recipient)
	msg.SetHeader("Subject", req.Title)
	// text/html for a html email
	msg.SetBody("text/plain", req.Body)

	n := gomail.NewDialer(s.conf.Host, s.conf.Port, s.conf.Username, s.conf.Username)
	// Send the email
	err := n.DialAndSend(msg)
	return err
}
