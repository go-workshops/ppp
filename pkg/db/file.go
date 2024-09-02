package db

import (
	"os"
)

type File struct {
	file *os.File
}

func OpenFile(name string) (*File, error) {
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	f := &File{
		file: file,
	}

	return f, nil
}

func (f *File) Write(bs []byte) (n int, err error) {
	return f.file.Write(bs)
}

func (f *File) Close() error {
	return f.file.Close()
}
