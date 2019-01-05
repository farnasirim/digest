package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/farnasirim/digest"
)

func TestSendEmail(t *testing.T) {
	var gmail digest.SMTPService = NewGmailSMTP("user@gmail.com", "pass")

	err := gmail.SendMailMultipart("receptionist@something.com", "hello",
		[]byte("body body body..."),
	)

	assert.NotNil(t, err)
}
