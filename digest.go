//go:generate mockgen --source=digest.go --destination=mock/mock.go --package=mock

package digest

type DocsService interface {
	TakeAndPersistSnapshopt(string) error
}

type SMTPService interface {
	SendMail(to string, msg []byte) error
	SendMailHtml(to, subject string, msg []byte) error
}

type DiffService interface {
	DiffDirsHtml(string, string) string
}
