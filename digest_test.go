package digest

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockGeneration(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "digest_go_mock")
	assert.Nil(t, err)
	defer func() {
		os.Remove(tmpFile.Name())
	}()

	fullMockgenString := fmt.Sprintf(`mockgen --source digest.go --destination %s --package=mock`, tmpFile.Name())
	mockgenCmdArgs := strings.Split(fullMockgenString, " ")
	mockgenCmd := exec.Command(mockgenCmdArgs[0], mockgenCmdArgs[1:]...)

	err = mockgenCmd.Start()
	assert.Nil(t, err)

	err = mockgenCmd.Wait()
	assert.Nil(t, err)

	originalMockContents, err := ioutil.ReadFile("./mock/mock.go")
	assert.Nil(t, err)

	newMockContents, err := ioutil.ReadFile(tmpFile.Name())
	assert.Nil(t, err)

	assert.Equal(t, newMockContents, originalMockContents)
}
