package store_test

import "anubot/store"

func Example() {
	store := store.New(store.HomePath())
	defer store.Close()
	store.InitDDL()
	// This is where you'd call API methods
}
