package drive

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteFile(t *testing.T) {
	dirName, err := ioutil.TempDir("", "some_prefix")
	assert.Nil(t, err)
	defer os.RemoveAll(dirName)

	filenames := []string{"f1.txt", "f2.html", "f3.bf"}
	contents := [][]byte{[]byte("somet content"), []byte("more"), []byte("!")}

	files := make([]*File, 0)
	for i := 0; i < len(filenames); i++ {
		files = append(files, NewFile(filenames[i], contents[i]))
	}

	driveService := NewDriveService(dirName, nil)
	snapshotName := "12-123-42"
	if err := driveService.writeSnapshot(snapshotName, files); err != nil {
		assert.Nil(t, err)
	}

	for i := 0; i < len(filenames); i++ {
		content, err := ioutil.ReadFile(path.Join(dirName, snapshotName, filenames[i]))
		assert.Nil(t, err)
		reflect.DeepEqual(content, contents[i])
	}

}
