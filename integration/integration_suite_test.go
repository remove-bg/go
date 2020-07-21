package integration_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var cliPath string

var _ = SynchronizedBeforeSuite(func() []byte {
	var err error
	path, err := gexec.Build("github.com/remove-bg/go")
	Expect(err).ShouldNot(HaveOccurred())
	return []byte(path)
}, func(data []byte) {
	cliPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {
	gexec.CleanupBuildArtifacts()
}, func() {})
