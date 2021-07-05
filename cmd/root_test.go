package cmd_test

import (
	. "./cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigureVersion", func() {
	It("sets the version", func() {
		ConfigureVersion("x.y.z", "sha")

		Expect(RootCmd.Version).To(Equal("x.y.z"))
	})
})
