package fileService

import "io"

type PhotosStorage interface {
	Read(path string) (io.ReadCloser, int64, error)
	Write(path string) (io.WriteCloser, error)
	Remove(path string) error
}
