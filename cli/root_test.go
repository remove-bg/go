package cli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/remove-bg/go/cli"
)

var _ = Describe("root", func() {
	It("has a version", func() {
		app := Bootstrap()
		Expect(app.Version).To(MatchRegexp(`\d+\.\d+\.\d+`))
	})
})
