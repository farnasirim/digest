package drive

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"google.golang.org/api/drive/v3"
)

var (
	ErrLessThanTwoSnapshots = errors.New("Less than two snapshots exist")
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

func (s *DriveService) TakeAndPersistTimedSnapshot(googleDocsFolder string) error {
	snapshotName := time.Now().Format(time.RFC3339)
	return s.TakeAndPersistSnapshot(snapshotName, googleDocsFolder)
}

func (s *DriveService) TakeAndPersistSnapshot(snapshotName, googleDocsFolder string) error {
	dirName := path.Join(s.dataDir, snapshotName)
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return err
	}

	folderQuery := `name = '%s' and mimeType = 'application/vnd.google-apps.folder'`
	folderQuery = fmt.Sprintf(folderQuery, googleDocsFolder)

	foldersResponse, err := s.driveSvc.Files.List().Q(folderQuery).Do()
	if err != nil {
		return err
	}

	if len(foldersResponse.Files) == 0 {
		log.Fatalf("No folder found with query: %s\n", folderQuery)
	} else if len(foldersResponse.Files) > 1 {
		log.Println("Multiple results returned... will use the first one!")
		for i, folder := range foldersResponse.Files {
			log.Println(i, folder.Id, folder.Name)
		}
	} else {
		log.Printf("Looking under folder %q with id %q",
			foldersResponse.Files[0].Name, foldersResponse.Files[0].Id)
	}

	docsFolder := foldersResponse.Files[0]

	filesUnderFolderQuery := `'%s' in parents`
	filesUnderFolderQuery = fmt.Sprintf(filesUnderFolderQuery, docsFolder.Id)

	filesUnderFolderResponse, err := s.driveSvc.Files.List().PageSize(500).
		Q(filesUnderFolderQuery).Fields("files(id, name, trashed)").Do()

	if err != nil {
		return err
	}

	for _, file := range filesUnderFolderResponse.Files {
		if file.Trashed {
			continue
		}
		resp, err := s.driveSvc.Files.Export(file.Id, "text/plain").Download()
		if err != nil {
			log.Printf("Error while processing %q, %q: %s\n", file.Name, file.Id, err.Error())
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		nameOnDisk := path.Join(dirName, file.Name+".txt")
		f, err := os.OpenFile(nameOnDisk, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Printf("Error while creating %q %q, %q: %s\n", nameOnDisk, file.Name, file.Id, err.Error())
			continue
		}
		if _, err := f.Write(body); err != nil {
			log.Printf("Error while writing %q, %q: %s\n", file.Name, file.Id, err.Error())
			continue
		}
		f.Close()

		log.Printf("Successfully written %q %q\n", file.Name, file.Id)
	}
	return nil
}

func (s *DriveService) writeSnapshot(snapshotName string, files []*File) error {
	pathPrefix := path.Join(s.dataDir, snapshotName)
	for _, file := range files {
		fullFilePath := path.Join(pathPrefix, file.relativePath)
		if err := os.MkdirAll(path.Dir(fullFilePath), 0755); err != nil {
			return err
		}
		if err := ioutil.WriteFile(fullFilePath, file.content, 0644); err != nil {
			return err
		}
	}
	return nil
}

func (s *DriveService) LastTwoDirs() (string, string, error) {
	dirs, err := ioutil.ReadDir(s.dataDir)
	if err != nil {
		return "", "", err
	}

	newest := ""
	for _, dir := range dirs {
		if dir.IsDir() && dir.Name() > newest {
			newest = dir.Name()
		}
	}

	beforeNewest := ""
	for _, dir := range dirs {
		if dir.IsDir() && dir.Name() > beforeNewest && dir.Name() != newest {
			beforeNewest = dir.Name()
		}
	}

	if beforeNewest == "" {
		return "", "", ErrLessThanTwoSnapshots
	}

	beforeNewest = path.Join(s.dataDir, beforeNewest)
	newest = path.Join(s.dataDir, newest)

	return beforeNewest, newest, nil
}
