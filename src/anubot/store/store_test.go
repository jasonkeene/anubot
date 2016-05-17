package store_test

import (
	"errors"
	"os/user"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "anubot/store"
)

var _ = Describe("Store", func() {
	var store *Store

	Describe("NewFromQuerier", func() {
		var mockQuerier *mockQuerier

		BeforeEach(func() {
			mockQuerier = newMockQuerier()
			store = NewFromQuerier(mockQuerier)
		})

		It("closes the querier when it is closed", func() {
			mockQuerier.CloseOutput.Ret0 <- nil
			store.Close()
			Expect(mockQuerier.CloseCalled).To(Receive())
		})

		Context("with a querier that errors on close", func() {
			It("returns that error", func() {
				err := errors.New("test-error")
				mockQuerier.CloseOutput.Ret0 <- err
				Expect(store.Close()).To(Equal(err))
			})
		})
	})
})

var _ = Describe("HomePath", func() {
	It("returns a file path in the user's home directory", func() {
		usr, _ := user.Current()
		Expect(HomePath()).To(Equal(path.Join(usr.HomeDir, "anubot.db")))
	})
})
