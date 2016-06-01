package bot_test

import (
	"strings"

	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

type equalLinesMatcher struct {
	matchers.EqualMatcher
}

func EqualLines(lines ...string) types.GomegaMatcher {
	matcher := &equalLinesMatcher{}
	matcher.Expected = []byte(strings.Join(lines, "\r\n") + "\r\n")
	return matcher
}
