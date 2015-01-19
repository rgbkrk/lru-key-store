package lrukeystore

import "testing"

func TestLRUKS(t *testing.T) {
	ks, err := New(2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if ks.IsIn("rgbkrk", "password") {
		t.Fatalf("Key in store on initialization\n")
	}

	ks.Add("rgbkrk", "password")

	if !ks.IsIn("rgbkrk", "password") {
		t.Fatalf("Key not in store\n")
	}

	// Fill the cache
	ks.Add("smashwilson", "myvoiceismypassport")
	if !ks.IsIn("smashwilson", "myvoiceismypassport") {
		t.Fatalf("New key not added and retrieved successfully")
	}
	ks.Add("ycombinator", "stillcantbelievethatsyourgithubusername")
	if !ks.IsIn("ycombinator", "stillcantbelievethatsyourgithubusername") {
		t.Fatalf("New key not added and retrieved successfully")
	}

	// The old rgbkrk key should be gone
	if ks.IsIn("rgbkrk", "password") {
		t.Fatalf("Key should have been cleared out of the cache, wasn't\ncache: %v\n", ks.cache)
	}

}
