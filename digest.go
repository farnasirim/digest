//go:generate mockgen --source=digest.go --destination=mock/mock.go --package=mock

package digest

import (
	"google.golang.org/api/drive/v3"
)

type GoogleAuthenticator interface {
	GetSvc() *drive.Service
}

type DocsService interface {
	TakeAndPersistSnapshopt(string) error
}

type SMTPService interface {
	SendMail(to string, msg []byte) error
	SendMailMultipart(to, subject string, msg []byte) error
}

type DiffService interface {
}
