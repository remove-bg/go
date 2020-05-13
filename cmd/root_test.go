package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/remove-bg/go/cmd"
)

var _ = Describe("CLI", func() {
	It("has a version", func() {
		Expect(RootCmd.Version).To(MatchRegexp(`\d+\.\d+\.\d+`))
	})
})
