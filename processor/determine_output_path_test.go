package processor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/remove-bg/go/processor"
)

var _ = Describe("DetermineOutputPath", func() {
	Context("when output path is set", func() {
		It("joins the original filename with the output path", func() {
			settings := Settings{
				OutputDirectory: "out",
			}

			result := DetermineOutputPath("in/nested/image.jpg", settings)

			Expect(result).To(Equal("out/image.png"))
		})
	})

	Context("when the output path isn't set", func() {
		It("writes to the original directory with a filename suffix", func() {
			settings := Settings{
				OutputDirectory: "",
			}

			result := DetermineOutputPath("in/nested/image.jpg", settings)

			Expect(result).To(Equal("in/nested/image-removebg.png"))
		})
	})

	Context("when the output format is set", func() {
		It("is used as the file extension", func() {
			settings := Settings{
				OutputDirectory: "out",
				ImageSettings: ImageSettings{
					OutputFormat: "jpg",
				},
			}

			result := DetermineOutputPath("in/nested/image.jpg", settings)

			Expect(result).To(Equal("out/image.jpg"))
		})
	})
})
