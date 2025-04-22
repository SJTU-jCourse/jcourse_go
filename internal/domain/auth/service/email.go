package service

import (
	"context"
	"fmt"

	"jcourse_go/pkg/mail"
)

type EmailService interface {
	SendVerificationCode(ctx context.Context, email string, code string) error
}

type emailService struct {
	sender       mail.EmailSender
	bodyTemplate string
	title        string
}

func (e *emailService) SendVerificationCode(ctx context.Context, email string, code string) error {
	body := fmt.Sprintf(e.bodyTemplate, code)
	err := e.sender.SendMail(ctx, email, e.title, body)
	if err != nil {
		return err
	}
	return err
}

func NewEmailService(bodyTemplate string, title string, sender mail.EmailSender) EmailService {
	return &emailService{
		bodyTemplate: bodyTemplate,
		title:        title,
		sender:       sender,
	}
}
