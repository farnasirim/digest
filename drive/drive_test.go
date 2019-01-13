package drive

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteSnapshot(t *testing.T) {
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

func TestLastTwoDirsNotEnough(t *testing.T) {
	dirName, err := ioutil.TempDir("", "some_prefix")
	assert.Nil(t, err)
	defer os.RemoveAll(dirName)

	err = os.Mkdir(path.Join(dirName, "something"), 0755)
	assert.Nil(t, err)

	driveService := NewDriveService(dirName, nil)
	_, _, err = driveService.LastTwoDirs()

	assert.Equal(t, ErrLessThanTwoSnapshots, err)
}

func TestLastTwoDirsNotEnoughWithFiles(t *testing.T) {
	dirName, err := ioutil.TempDir("", "some_prefix")
	assert.Nil(t, err)
	defer os.RemoveAll(dirName)

	err = os.Mkdir(path.Join(dirName, "something"), 0755)
	assert.Nil(t, err)

	f, err := os.OpenFile(path.Join(dirName, "a"),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	assert.Nil(t, err)
	defer f.Close()

	dirs, err := ioutil.ReadDir(dirName)
	assert.Nil(t, err)
	assert.True(t, 2 <= len(dirs))

	driveService := NewDriveService(dirName, nil)
	_, _, err = driveService.LastTwoDirs()

	assert.Equal(t, ErrLessThanTwoSnapshots, err)
}

func TestLastTwoDirsSimple(t *testing.T) {
	dirName, err := ioutil.TempDir("", "some_prefix")
	assert.Nil(t, err)
	defer os.RemoveAll(dirName)

	f, err := os.OpenFile(path.Join(dirName, "a"),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	assert.Nil(t, err)
	defer f.Close()

	f, err = os.OpenFile(path.Join(dirName, "z"),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	assert.Nil(t, err)
	defer f.Close()

	err = os.Mkdir(path.Join(dirName, "b"), 0755)
	assert.Nil(t, err)

	err = os.Mkdir(path.Join(dirName, "c"), 0755)
	assert.Nil(t, err)

	dirs, err := ioutil.ReadDir(dirName)
	assert.Nil(t, err)
	assert.True(t, 2 <= len(dirs))

	driveService := NewDriveService(dirName, nil)
	older, newer, err := driveService.LastTwoDirs()

	assert.Nil(t, err)
	assert.Equal(t, older, "b")
	assert.Equal(t, newer, "c")
}
