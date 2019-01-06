package drive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthGoogle(t *testing.T) {
	gAuth := NewGoogleAuthenticator("")
	svc, err := gAuth.GetOrCreateClient()
	assert.True(t, (err == nil) != (svc == nil))
}
