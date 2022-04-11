package impl

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"fmt"
	"io"
	"os"
)

const testCreateFilename = "test_create_filename"

type fileServicePhotos struct {
	root string
}

func NewFileServicePhotos(root string) (*fileServicePhotos, error) {
	if root[len(root)-1] != '/' {
		root += "/"
	}
	//check creating and deleting files perms
	openedFile, err := os.Open(root + testCreateFilename)
	if err == nil {
		err = openedFile.Close()
		if err != nil {
			return nil, fmt.Errorf("close opened test file failed: %w", err)
		}
		err = os.Remove(root + testCreateFilename)
		if err != nil {
			return nil, fmt.Errorf("removing test file failed: %w", err)
		}
	}

	createdFile, err := os.Create(root + testCreateFilename)
	if err != nil {
		return nil, fmt.Errorf("creating test file failed: %w", err)
	}
	err = createdFile.Close()
	if err != nil {
		return nil, fmt.Errorf("close opened test file failed: %w", err)
	}
	err = os.Remove(root + testCreateFilename)
	if err != nil {
		return nil, fmt.Errorf("removing test file failed: %w", err)
	}

	return &fileServicePhotos{root: root}, nil
}

func (repo *fileServicePhotos) Read(path string) (io.ReadCloser, int64, error) {
	openedFile, err := os.Open(repo.root + path)
	if err != nil {
		return nil, 0, handlers.ErrBaseApp.Wrap(err, "read opening file failed")
	}
	stats, err := openedFile.Stat()
	if err != nil {
		return nil, 0, handlers.ErrBaseApp.Wrap(err, "get stats file failed")
	}
	return openedFile, stats.Size(), nil
}

func (repo *fileServicePhotos) Write(path string) (io.WriteCloser, error) {
	openedFile, err := os.Create(repo.root + path)
	if err != nil {
		return nil, handlers.ErrBaseApp.Wrap(err, "write creating file failed")
	}
	return openedFile, nil
}

func (repo *fileServicePhotos) Remove(path string) error {
	err := os.Remove(path)
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "removing file failed")
	}
	return nil
}
