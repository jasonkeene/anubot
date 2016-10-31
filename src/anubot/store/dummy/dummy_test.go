package dummy

import (
	"testing"

	"anubot/store"
)

func TestThatDummyBackendCompliesWithAllStoreMethods(t *testing.T) {
	var _ store.Store = &Dummy{}
}
