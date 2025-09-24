package email

import "context"

type Request struct {
	Title     string
	Body      string
	Recipient string
}

type EmailSender interface {
	SendEmail(ctx context.Context, req Request) error
}
