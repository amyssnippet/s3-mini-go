package storage

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteStream(t *testing.T) {
    tmpDir := t.TempDir()

	var keyPath string

	masterKey, err := loadOrGenerateStorageKey(keyPath)
	if err != nil {
		log.Fatalf("Failed to load storage key: %v", err)
	}
    store := NewStore(tmpDir, masterKey)

    data := []byte("Hello S3 Mini")
    reader := bytes.NewReader(data)

    n, err := store.WriteStream("test_file.txt", reader)
    if err != nil {
        t.Fatalf("Failed to write: %v", err)
    }

    if n != int64(len(data)) {
        t.Errorf("Expected %d bytes, wrote %d", len(data), n)
    }

    finalPath := filepath.Join(tmpDir, "test_file.txt")
    if _, err := os.Stat(finalPath); os.IsNotExist(err) {
        t.Error("File was not created on disk")
    }
}

func loadOrGenerateStorageKey(path string) ([]byte, error) {
	keyPath := filepath.Join(path, "storage.key")
	key, err := os.ReadFile(keyPath)
	if err == nil {
		if len(key) != 32 {
			return nil, fmt.Errorf("invalid key length in %s", keyPath)
		}
		return key, nil
	}
	fmt.Println("Generating new storage encryption key...")
	newKey := make([]byte, 32)
	if _, err := rand.Read(newKey); err != nil {
		return nil, err
	}
	if err := os.WriteFile(keyPath, newKey, 0600); err != nil {
		return nil, err
	}
	return newKey, nil
}