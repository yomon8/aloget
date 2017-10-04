package fileio

import (
	"io"
	"os"
)

type FileIO interface {
	GetReadStream(*os.File) (*io.Reader, error)
	GetWriteFile(*string) (*os.File, error)
}
