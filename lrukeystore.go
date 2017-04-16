package lrukeystore

import (
	"crypto/rand"
	"errors"

	"github.com/hashicorp/golang-lru"
	"golang.org/x/crypto/bcrypt"
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
	val, ok := ks.cache.Get(user)
	if !ok {
		// Fake a bcrypt comparison
		bcrypt.CompareHashAndPassword([]byte(""), []byte(putative))
		return false
	}

	expectedHash := val.([]byte)

	return bcrypt.CompareHashAndPassword(expectedHash, []byte(putative)) == nil
}

// Add adds a new key to the KeyStore using the internal hashing scheme
func (ks *KeyStore) Add(user string, key string) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)

	// TODO: Change tests and interface to return error
	if err != nil {
		panic(err)
	}

	ks.cache.Add(user, hashed)
}
