package store_test

import (
	"math/rand"
	"path"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "anubot/store"
)

var _ = Describe("Store API Methods", func() {
	var store *Store

	BeforeEach(func() {
		store = New(path.Join(tempDir, "test-"+strconv.Itoa(rand.Int())))
		Expect(store.InitDDL()).To(Succeed())
	})

	AfterEach(func() {
		store.Close()
	})

	It("stores user credentials", func() {
		expectedUser, expectedPass := "test-user", "test-pass"
		kind := "test-kind"

		_, _, err := store.Credentials(kind)
		Expect(err).To(HaveOccurred())

		err = store.SetCredentials(kind, expectedUser, expectedPass)
		Expect(err).ToNot(HaveOccurred())

		username, pass, err := store.Credentials(kind)
		Expect(err).ToNot(HaveOccurred())
		Expect(username).To(Equal(expectedUser))
		Expect(pass).To(Equal(expectedPass))
	})

	It("can tell if it has valid user credentials or not", func() {
		expectedUser, expectedPass := "test-user", "test-pass"
		kind := "test-kind"

		Expect(store.HasCredentials(kind)).To(BeFalse())

		err := store.SetCredentials(kind, expectedUser, expectedPass)
		Expect(err).ToNot(HaveOccurred())

		Expect(store.HasCredentials(kind)).To(BeTrue())

		err = store.SetCredentials(kind, "", expectedPass)
		Expect(err).ToNot(HaveOccurred())

		Expect(store.HasCredentials(kind)).To(BeFalse())

		err = store.SetCredentials(kind, expectedUser, "")
		Expect(err).ToNot(HaveOccurred())

		Expect(store.HasCredentials(kind)).To(BeFalse())
	})
})
