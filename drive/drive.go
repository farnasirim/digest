package drive

import (
	"io/ioutil"
	"os"
	"path"

	"google.golang.org/api/drive/v3"
)

type DriveService struct {
	dataDir  string
	driveSvc *drive.Service
}

type File struct {
	relativePath string
	content      []byte
}

func NewDriveService(dataDir string, driveSvc *drive.Service) *DriveService {
	ret := &DriveService{
		dataDir:  dataDir,
		driveSvc: driveSvc,
	}

	return ret
}

func NewFile(relativePath string, content []byte) *File {
	return &File{
		relativePath: relativePath,
		content:      content,
	}
}

func (s *DriveService) TakeAndPersistSnapshot(snapshotName string) error {
	return nil
}

func (s *DriveService) writeSnapshot(snapshotName string, files []*File) error {
	pathPrefix := path.Join(s.dataDir, snapshotName)
	for _, file := range files {
		fullFilePath := path.Join(pathPrefix, file.relativePath)
		println(" :: ", fullFilePath)
		println(" @@ ", path.Dir(fullFilePath))
		if err := os.MkdirAll(path.Dir(fullFilePath), os.ModePerm); err != nil {
			return err
		}
		if err := ioutil.WriteFile(fullFilePath, file.content, 0644); err != nil {
			return err
		}
	}
	return nil
}
