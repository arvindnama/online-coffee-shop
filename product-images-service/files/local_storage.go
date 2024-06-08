package files

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

// / LocalStorage is an local on disk storage implementation
type LocalStorage struct {
	maxFileSize int
	basePath    string
}

func NewLocalStorage(basePath string, maxFileSize int) (*LocalStorage, error) {
	absPath, err := filepath.Abs(basePath)

	if err != nil {

		return nil, err

	}
	return &LocalStorage{
		basePath:    absPath,
		maxFileSize: maxFileSize,
	}, nil
}

func (l *LocalStorage) Save(path string, file io.Reader) error {
	fp := l.fullPath(path)

	dir := filepath.Dir(fp)

	err := os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return xerrors.Errorf("Unable to create directory %w", err)
	}

	_, err = os.Stat(fp)

	if os.IsExist(err) {
		err = os.Remove(fp)
		if err != nil {
			return xerrors.Errorf("Unable to delete file %w", err)
		}
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf("Unable to get file info %w", err)
	}

	// create file
	f, err := os.Create(fp)

	if err != nil {
		return xerrors.Errorf("Unable to create file %w", err)
	}

	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return xerrors.Errorf("Unable to write to file %w", err)
	}
	return nil
}

func (l *LocalStorage) Get(path string) (*os.File, error) {
	fp := l.fullPath(path)

	f, err := os.Open(fp)

	if err != nil {
		return nil, xerrors.Errorf("Unable to open file %w", err)
	}
	return f, nil
}

func (l *LocalStorage) fullPath(path string) string {
	return filepath.Join(l.basePath, path)
}
