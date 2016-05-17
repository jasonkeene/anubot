package bot_test

import (
	"strings"

	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

type equalLines struct {
	matchers.EqualMatcher
}

func EqualLines(lines ...string) types.GomegaMatcher {
	expected := []byte(strings.Join(lines, "\r\n") + "\r\n")
	matcher := &equalLines{}
	matcher.Expected = expected
	return matcher
}
