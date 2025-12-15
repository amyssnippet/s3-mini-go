package security

import (
    "testing"
    "os"
)

func TestKeyStore(t *testing.T) {
    tmpFile := "test_keys.json"
    defer os.Remove(tmpFile)

    ks, _ := NewKeyStore(tmpFile)
    
    err := ks.CreateKey("sk_test_123", "WRITE")
    if err != nil {
        t.Errorf("Failed to create key: %v", err)
    }

    if !ks.IsAllowed("sk_test_123", "WRITE") {
        t.Error("Key should have WRITE permission")
    }
    
    if ks.IsAllowed("sk_test_123", "ADMIN") {
        t.Error("Key should NOT have ADMIN permission")
    }
}