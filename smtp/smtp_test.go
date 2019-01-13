package smtp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/farnasirim/digest"
)

func TestSendEmail(t *testing.T) {
	var gmail digest.SMTPService = NewSimpleSMTP("smtp.gmail.com", "587",
		"user@gmail.com", "pass")
	err := gmail.SendMailHtml("receptionist@something.com", "hello",
		[]byte("body body body..."),
	)

	assert.NotNil(t, err)
}
