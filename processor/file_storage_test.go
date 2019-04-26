package processor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"path"
	"runtime"

	. "github.com/remove-bg/go/processor"
)

var _ = Describe("FileStorage", func() {
	var (
		subject FileStorage
	)

	BeforeEach(func() {
		subject = FileStorage{}
	})

	Describe("FileExists", func() {
		It("is true when the file is present", func() {
			_, testFile, _, _ := runtime.Caller(0)
			fixtureFile := path.Join(path.Dir(testFile), "../fixtures/person-in-field.jpg")

			Expect(subject.FileExists(fixtureFile)).To(BeTrue())
		})

		It("is false when the file doesn't exist", func() {
			_, testFile, _, _ := runtime.Caller(0)
			missing := path.Join(path.Dir(testFile), "../fixtures/missing.jpg")

			Expect(subject.FileExists(missing)).To(BeFalse())
		})

		It("is false when the directory doesn't exist", func() {
			_, testFile, _, _ := runtime.Caller(0)
			missingDir := path.Join(path.Dir(testFile), "../missing")

			Expect(subject.FileExists(missingDir)).To(BeFalse())
		})
	})
})
