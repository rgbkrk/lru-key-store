# lru-key-store :recycle: :key:

Least Recently Used Key Store

![RAM Keychain Cache](https://cloud.githubusercontent.com/assets/836375/5823789/40fe31e0-a0a5-11e4-8486-1ecffe0dc86c.jpg)

Cache secrets (API keys, passwords) per user without storing them in plaintext, rotating out the least recently used.

Entries are rotated out using the [Least Recently Used](http://en.wikipedia.org/wiki/Cache_algorithms#LRU) cache algorithm. Under the hood, we're using [hashicorp/golang-lru](https://github.com/hashicorp/golang-lru) as a thread safe fixed-size cache and mapping `user -> HMAC-SHA(user secret)`. The secret key for the HMAC-SHA is initialized per LRU Key Store.

## High level usage

LRU Key Store was created to be a caching layer between itself and a source of truth. If a key isn't in the store, we test the putative API key provided by the user with the upstream source of truth. If that passes, we add it to the cache. If it doesn't, nothing changes. If a key is in the store, we just respond back immediately. Invalid keys should then take the same amount of time to repudiate while valid keys are instantaneous.

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
	// The key store should not have this password
	in := ks.IsIn("rgbkrk", password)
	fmt.Printf("'%v' -> %v\n", password, in)

	// Nor another key
	password = "secret"
	fmt.Printf("'%v' -> %v\n", password, ks.IsIn("rgbkrk", password))

	// Give the keystore a secret
	ks.Add("rgbkrk", "secret")

	// Invalid key
	password = "password"
	fmt.Printf("'%v' -> %v\n", password, ks.IsIn("rgbkrk", password))

	// We have this key!
	password = "secret"
	fmt.Printf("'%v' -> %v\n", password, ks.IsIn("rgbkrk", password))

	// Fill the cache
	ks.Add("smashwilson", "myvoiceismypassport")
	ks.Add("ycombinator", "stillcantbelievethatsyourgithubusername")

	// rgbkrk should no longer have a key stored
	password = "secret"
	fmt.Printf("'%v' -> %v\n", password, ks.IsIn("rgbkrk", password))

}
```
