package processor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"github.com/sirupsen/logrus/hooks/test"

	. "github.com/remove-bg/go/processor"
)

var _ = Describe("Notifier", func() {
	Describe("Success", func() {
		It("logs the image details", func() {
			logger, hook := test.NewNullLogger()
			subject := Notifier{
				Logger: logger,
			}

			subject.Success("input/image.jpg", 1, 2)

			logged := hook.LastEntry()

			Expect(logged).ToNot(BeNil())
			Expect(logged.Message).To(Equal("Processed image"))
			Expect(logged.Data["image"]).To(Equal("1/2"))
			Expect(logged.Data["input"]).To(Equal("input/image.jpg"))
		})
	})

	Describe("Error", func() {
		It("logs the error and image details", func() {
			logger, hook := test.NewNullLogger()
			subject := Notifier{
				Logger: logger,
			}

			err := errors.New("boom")
			subject.Error(err, "input/image.jpg", 1, 2)

			logged := hook.LastEntry()

			Expect(logged).ToNot(BeNil())
			Expect(logged.Message).To(Equal("boom"))
			Expect(logged.Data["image"]).To(Equal("1/2"))
			Expect(logged.Data["input"]).To(Equal("input/image.jpg"))
		})
	})

	Describe("NewNotifier", func() {
		It("builds a notifier", func() {
			n := NewNotifier()
			Expect(n.Logger).ToNot(BeNil())
		})
	})
})
