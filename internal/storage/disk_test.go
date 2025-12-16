package storage

import (
    "os"
    "testing"
    "bytes"
    "path/filepath"
)

func TestWriteStream(t *testing.T) {
    tmpDir := t.TempDir()
    store := NewStore(tmpDir)

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