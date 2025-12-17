package core

import "testing"

func TestFileMetadata(t *testing.T) {
    meta := FileMetadata{
        Name: "test.txt",
        Size: 1024,
    }
    if meta.Name != "test.txt" {
        t.Error("Metadata struct is malformed")
    }
}