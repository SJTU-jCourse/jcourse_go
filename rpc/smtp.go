package rpc

import (
	"context"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
	"jcourse_go/util"
)

func SendMail(ctx context.Context, recipient string, subject string, body string) error {
	if util.IsDebug() {
		return nil
	}
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	smtpSender := os.Getenv("SMTP_SENDER")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	msg := gomail.NewMessage()
	msg.SetHeader("From", smtpSender)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", subject)
	// text/html for a html email
	msg.SetBody("text/plain", body)

	n := gomail.NewDialer(host, int(port), username, password)

	// Send the email
	err = n.DialAndSend(msg)
	return err
}
