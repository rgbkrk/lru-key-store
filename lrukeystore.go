package lrukeystore

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"

	"github.com/hashicorp/golang-lru"
)

// KeyStore is a fixed size cache of keys using an LRU cache
type KeyStore struct {
	cache     *lru.Cache
	systemKey []byte
}

// New creates a KeyStore of the given size
func New(size int) (*KeyStore, error) {
	cache, err := lru.New(size)

	if err != nil {
		return nil, err
	}

	systemKeySize := 64
	systemKey := make([]byte, systemKeySize)

	n, err := rand.Read(systemKey)

	if err != nil {
		return nil, err
	}

	if n != systemKeySize {
		return nil, errors.New("Unable to allocate system key")
	}

	keyStore := &KeyStore{
		cache:     cache,
		systemKey: systemKey,
	}
	return keyStore, nil

}

// IsIn checks to see if user has the putative key in the KeyStore
func (ks *KeyStore) IsIn(user string, putative string) bool {
	mac := hmac.New(sha256.New, ks.systemKey)
	mac.Write([]byte(putative))
	computedMAC := mac.Sum(nil)

	val, ok := ks.cache.Get(user)
	if !ok {
		return false
	}

	expectedMAC := val.([]byte)

	return hmac.Equal(expectedMAC, computedMAC)
}

// Add adds a new key to the KeyStore using the internal hashing scheme
func (ks *KeyStore) Add(user string, key string) {
	mac := hmac.New(sha256.New, ks.systemKey)
	mac.Write([]byte(key))
	expectedMAC := mac.Sum(nil)

	ks.cache.Add(user, expectedMAC)
}
