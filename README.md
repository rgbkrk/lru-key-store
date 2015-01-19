# lru-key-store
:recycle::key: Least Recently Used Key Store

### Example usage

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

When a key is unavailable, it's the implementor's decision to request a new key from a source of truth. I'd like to use this for caching API keys before hitting some service that has rate limits.




