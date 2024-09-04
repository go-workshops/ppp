package db

import (
	"os"
)

type FS struct {
	dir  string
	file string
}

func OpenFS(name string) (FS, error) {
	err := os.MkdirAll(name, 0755)
	if err != nil {
		return FS{}, err
	}

	return FS{dir: name}, nil
}

func (fs FS) File(name string) FS {
	return FS{
		dir:  fs.dir,
		file: name,
	}
}

func (fs FS) Write(data []byte) (n int, err error) {
	if fs.file == "" {
		return 0, nil
	}

	err = os.WriteFile(fs.dir+"/"+fs.file, data, 0644)
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

func (fs FS) Read(bs []byte) (n int, err error) {
	if fs.file == "" {
		return 0, os.ErrInvalid
	}

	data, err := os.ReadFile(fs.dir + "/" + fs.file)
	if err != nil {
		return 0, err
	}

	return copy(bs, data), nil
}
