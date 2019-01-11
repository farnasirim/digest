package smtp

import (
	"fmt"
	"net/smtp"
)

type SimpleSMTP struct {
	fromAddress string
	password    string
}

func NewSimpleSMTP(fromAddress, password string) *SimpleSMTP {
	ret := &SimpleSMTP{
		fromAddress: fromAddress,
		password:    password,
	}

	return ret
}

func (s *SimpleSMTP) SendMail(to string, msg []byte) error {
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", s.fromAddress, s.password, "smtp.gmail.com"),
		s.fromAddress, []string{to}, msg,
	)

	return err
}

func (s *SimpleSMTP) SendMailMultipart(to, subject string, msg []byte) error {
	completedMessage := append([]byte(
		"From: "+s.fromAddress+"\n"+
			"To: "+to+"\n"+
			`Content-Type: text/html; charset="utf-8"`+"\n"+
			fmt.Sprintf("Subject: %s\n\n", subject),
	), msg...)

	return s.SendMail(to, completedMessage)
}
