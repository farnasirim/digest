package drive

import (
	_ "golang.org/x/net/context"
	_ "golang.org/x/oauth2"
	_ "golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type GoogleAuthenticator struct {
}

func NewGoogleAuthenticator() *GoogleAuthenticator {
	auth := &GoogleAuthenticator{}

	return auth
}
func (a *GoogleAuthenticator) GetSvc() (*drive.Service, error) {
	return nil, nil
}
