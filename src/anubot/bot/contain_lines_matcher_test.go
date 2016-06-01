package bot_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type containLinesMatcher struct {
	MatchBytes []byte
}

func ContainLines(lines ...string) types.GomegaMatcher {
	matcher := &containLinesMatcher{}
	matcher.MatchBytes = []byte(strings.Join(lines, "\r\n") + "\r\n")
	return matcher
}

func (m *containLinesMatcher) Match(actual interface{}) (success bool, err error) {
	actualBytes, ok := actual.([]byte)
	if !ok {
		return false, fmt.Errorf("ContainLines matcher requires a []byte.  Got:\n%s", format.Object(actual, 1))
	}

	return bytes.Contains(actualBytes, m.MatchBytes), nil
}

func (m *containLinesMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to contain bytes", m.MatchBytes)
}

func (m *containLinesMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to contain bytes", m.MatchBytes)
}
