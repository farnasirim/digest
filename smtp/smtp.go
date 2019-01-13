package smtp

import (
	"fmt"
	"net/smtp"
)

type SimpleSMTP struct {
	host        string
	port        string
	fromAddress string
	password    string
}

func NewSimpleSMTP(host, port, fromAddress, password string) *SimpleSMTP {
	ret := &SimpleSMTP{
		host:        host,
		port:        port,
		fromAddress: fromAddress,
		password:    password,
	}

	return ret
}

func (s *SimpleSMTP) SendMail(to string, msg []byte) error {
	err := smtp.SendMail(s.host+":"+s.port,
		smtp.PlainAuth("", s.fromAddress, s.password, s.host),
		s.fromAddress, []string{to}, msg,
	)

	return err
}

func (s *SimpleSMTP) SendMailHtml(to, subject string, msg []byte) error {
	completedMessage := append([]byte(
		"From: "+s.fromAddress+"\n"+
			"To: "+to+"\n"+
			`Content-Type: text/html; charset="utf-8"`+"\n"+
			fmt.Sprintf("Subject: %s\n\n", subject),
	), msg...)

	return s.SendMail(to, completedMessage)
}
