package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"s3-mini/internal/security"
)

type Store struct {
	Root      string
	MasterKey []byte
}

func NewStore(root string, key []byte) *Store {
	os.MkdirAll(root, 0755)
	return &Store{
		Root:      root,
		MasterKey: key,
	}
}

func (s *Store) WriteStream(name string, r io.Reader) (int64, error) {
	path := filepath.Join(s.Root, filepath.Base(name))
	f, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	secureWriter, err := security.NewEncryptedWriter(f, s.MasterKey)
	if err != nil {
		return 0, fmt.Errorf("encryption init failed: %v", err)
	}
	return io.Copy(secureWriter, r)
}

func (s *Store) ReadStream(name string) (io.ReadCloser, error) {
    path := filepath.Join(s.Root, filepath.Base(name))
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    secureReader, err := security.NewDecryptedReader(f, s.MasterKey)
    if err != nil {
        f.Close()
        return nil, fmt.Errorf("decryption init failed: %v", err)
    }
    return &wrappedReadCloser{Reader: secureReader, Closer: f}, nil
}


type wrappedReadCloser struct {
    io.Reader
    io.Closer
}