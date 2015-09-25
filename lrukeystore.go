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

func (ks *KeyStore) deriveHash(user, key string) []byte {
	mac := hmac.New(sha256.New, ks.systemKey)
	mac.Write([]byte(user))
	mac.Write([]byte(key))
	return mac.Sum(nil)
}

// IsIn checks to see if user has the putative key in the KeyStore
func (ks *KeyStore) IsIn(user, putative string) bool {

	computedMAC := ks.deriveHash(user, putative)

	_, ok := ks.cache.Get(string(computedMAC))

	return ok
}

// Add adds a new key to the KeyStore using the internal hashing scheme
func (ks *KeyStore) Add(user, key string) {
	expectedMAC := ks.deriveHash(user, key)

	ks.cache.Add(string(expectedMAC), true)
}
