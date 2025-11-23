package storage

import (
	"io"
	"os"
	"path/filepath"
)


type Store struct {
	Root string
}

func NewStore(root string) *Store{
	os.MkdirAll(root, 0755)
	return &Store{Root: root}
}

func (s *Store) WriteStream(filename string, r io.Reader) (int64, error) {
	finalPath := filepath.Join(s.Root, filepath.Base(filename))

	f, err := os.Create(finalPath)

	if err != nil {
		return 0, err
	}

	defer f.Close()

	return io.Copy(f, r)
}