package main

import (
	"fmt"
	"net/smtp"
)

type GmailSMTP struct {
	fromAddress string
	password    string
}

func NewGmailSMTP(fromAddress, password string) *GmailSMTP {
	ret := &GmailSMTP{
		fromAddress: fromAddress,
		password:    password,
	}

	return ret
}

func (s *GmailSMTP) SendMail(to string, msg []byte) error {
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", s.fromAddress, s.password, "smtp.gmail.com"),
		s.fromAddress, []string{to}, msg,
	)

	return err
}

func (s *GmailSMTP) SendMailMultipart(to, subject string, msg []byte) error {
	completedMessage := append([]byte(
		"From: "+s.fromAddress+"\n"+
			"To: "+to+"\n"+
			`Content-Type: text/html; charset="utf-8"`+"\n"+
			fmt.Sprintf("Subject: %s\n\n", subject),
	), msg...)

	return s.SendMail(to, completedMessage)
}
