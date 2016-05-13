package bot_test

import (
	"strings"

	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

type equalLines struct {
	expected []byte
	matchers.EqualMatcher
}

func EqualLines(lines ...string) types.GomegaMatcher {
	expected := []byte(strings.Join(lines, "\r\n") + "\r\n")
	return &equalLines{
		expected:     expected,
		EqualMatcher: matchers.EqualMatcher{Expected: expected},
	}
}
