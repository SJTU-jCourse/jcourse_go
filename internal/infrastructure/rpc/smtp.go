package rpc

import (
	"context"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"

	"jcourse_go/pkg/util"
)

func SendMail(ctx context.Context, recipient string, subject string, body string) error {
	if util.IsDebug() {
		return nil
	}
	host := util.GetSMTPHost()
	portStr := util.GetSMTPPort()
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	smtpSender := util.GetSMTPSender()
	username := util.GetSMTPUser()
	password := util.GetSMTPPassword()

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
