package security

import (
	"encoding/json"
	"os"
	"sync"
)


const (
	PermRead = "READ"
	PermWrite = "WRITE"
	PermAdmin = "ADMIN"
)


type KeyStore struct {
	mu sync.RWMutex
	Keys map[string]string
	path string
}

func NewKeyStore(path string) (*KeyStore, error) {
	ks := &KeyStore{
		Keys: make(map[string]string),
		path: path,
	}

	if data, err := os.ReadFile(path); err == nil {
		json.Unmarshal(data, &ks.Keys)
	}

	return ks, nil
}


func (ks *KeyStore) IsAllowed(apiKey string, requiredPerm string) bool {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	role, exists := ks.Keys[apiKey]

	if !exists {
		return false
	}

	if role == PermAdmin {
		return true
	}

	return role == requiredPerm
}


func (ks *KeyStore) CreateKey(apiKey string, role string) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	ks.Keys[apiKey] = role

	return ks.save()
}

func (ks *KeyStore) save() error {
	data, err := json.MarshalIndent(ks.Keys, "", "	")

	if err != nil {
		return err
	}

	return os.WriteFile(ks.path, data, 0600)
}