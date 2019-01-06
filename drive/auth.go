package drive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type GoogleAuthenticator struct {
	secretDir string
}

const (
	credentialsFileName = "credentials.json"
	tokenFileName       = "token.json"
)

func NewGoogleAuthenticator(secretDir string) *GoogleAuthenticator {
	auth := &GoogleAuthenticator{
		secretDir: secretDir,
	}

	return auth
}

// tokenFromFile Retrieves a token from a local file.
func (a *GoogleAuthenticator) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func (a *GoogleAuthenticator) tokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

// tokenToFile saves an oauth2 token to a file.
func (a *GoogleAuthenticator) tokenToFile(token *oauth2.Token, filePath string) {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to save oauth token to %s: (%v)", filePath, err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// GetOrCreateClient reads credentials.json from the supplied secretDir.
// It then goes on to either create a client object from the saved token.json
// or create one and persist the token.json in the sercertDir.
func (a *GoogleAuthenticator) GetOrCreateClient() (*http.Client, error) {
	credentials, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(credentials, drive.DriveScope)
	if err != nil {
		return nil, err
	}

	tok, err := a.tokenFromFile(path.Join(a.secretDir, tokenFileName))
	if err != nil {
		log.Println("Could not get token from file. Trying web...")

		tok, err = a.tokenFromWeb(config)
		if err != nil {
			return nil, err
		}

		a.tokenToFile(tok, path.Join(a.secretDir, tokenFileName))
	}

	client := config.Client(context.Background(), tok)

	return client, nil
}
