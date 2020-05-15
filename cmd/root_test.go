package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/remove-bg/go/cmd"
)

var _ = Describe("ConfigureVersion", func() {
	It("sets the version", func() {
		ConfigureVersion("x.y.z", "sha")

		Expect(RootCmd.Version).To(Equal("x.y.z"))
	})
})
