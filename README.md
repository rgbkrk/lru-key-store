# lru-key-store
:recycle::key: Least Recently Used Key Store

Cache keys (API keys, passwords) per user without storing them in plaintext, rotating out the least recently used.

Entries are rotated out using the [Least Recently Used](http://en.wikipedia.org/wiki/Cache_algorithms#LRU) cache algorithm, implemented by [hashicorp/golang-lru](https://github.com/hashicorp/golang-lru). The cache is a fixed size thread safe LRU cache

This was created to cache API keys without storing the API keys in plaintext. Internally, this currently maps `user -> HMAC-SHA(api-key)`. The secret key is generated per initialization of an LRU key store.

## Features

* Add key
* Check if key is currently cached

## Example usage

```go
package main

import (
	"fmt"

	lruks "github.com/rgbkrk/lru-key-store"
)

func main() {
	// Cache a whopping 2 keys
	ks, err := lruks.New(2)
	if err != nil {
		panic(err)
	}

	password := "password"
	// Does the keystore have this password?
	in := ks.IsIn("rgbkrk", password)
	fmt.Printf("'%v' -> %v\n", password, in)

	// How about this one?
	password = "secret"
	fmt.Printf("'%v' -> %v\n", password, ks.IsIn("rgbkrk", password))

	// Give the keystore a secret
	ks.Add("rgbkrk", "secret")

	password = "password"
	fmt.Printf("'%v' -> %v\n", password, ks.IsIn("rgbkrk", password))

	password = "secret"
	fmt.Printf("'%v' -> %v\n", password, ks.IsIn("rgbkrk", password))

	// Fill the cache
	ks.Add("smashwilson", "myvoiceismypassport")
	ks.Add("ycombinator", "stillcantbelievethatsyourgithubusername")

	// Check the key for rgbkrk
	password = "secret"
	fmt.Printf("'%v' -> %v\n", password, ks.IsIn("rgbkrk", password))

}
```

## Design Goal

LRU Key Store was created to be a caching layer between itself and a source of truth. If a key isn't in the store, we use the key provided with the upstream source of truth. If that passes, we add it to the cache. If it doesn't, nothing changes. If a key is in the store, we just respond back immediately.
