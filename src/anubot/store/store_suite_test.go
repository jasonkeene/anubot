package store_test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store Suite")
}

var tempDir string

var _ = BeforeSuite(func() {
	// create temp dir to store databases for each test
	var err error
	tempDir, err = ioutil.TempDir("", "store-test-suite")
	Expect(err).ToNot(HaveOccurred())

	// seed RNG
	rand.Seed(time.Now().UnixNano())
})

var _ = AfterSuite(func() {
	// clean up tmp dir
	os.RemoveAll(tempDir)
})
