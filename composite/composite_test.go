package composite_test

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"./composite"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
)

var _ = Describe("Composite", func() {
	var (
		subject    composite.Compositor
		exampleZip string
		outputPath string
		testDir    string
	)

	BeforeEach(func() {
		subject = composite.New()

		_, testFile, _, _ := runtime.Caller(0)
		testDir = path.Dir(testFile)

		exampleZip = path.Join(testDir, "../fixtures/zip/example-cat.zip")
		outputPath = path.Join(testDir, fmt.Sprintf("../tmp/composite-cat-%d.png", config.GinkgoConfig.ParallelNode))

		// Remove stale state from any previous test runs
		os.Remove(outputPath)
	})

	Context("when the input zip does not exist", func() {
		It("returns an error", func() {
			Expect(subject.Process("missing.zip", outputPath)).To(MatchError("Could not locate zip: missing.zip"))
		})

		It("does not write any output", func() {
			Expect(subject.Process("missing.zip", outputPath)).To(HaveOccurred())
			Expect(outputPath).ToNot(BeAnExistingFile())
		})
	})

	Context("when color.jpg does not exist in the input zip", func() {
		It("returns an error", func() {
			exampleZip = path.Join(testDir, "../fixtures/zip/example-missing-color.zip")
			Expect(exampleZip).To(BeAnExistingFile())

			Expect(subject.Process(exampleZip, outputPath)).To(MatchError("Unable to find image in ZIP: color.jpg"))
		})
	})

	Context("when alpha.png does not exist in the input zip", func() {
		It("returns an error", func() {
			exampleZip = path.Join(testDir, "../fixtures/zip/example-missing-alpha.zip")
			Expect(exampleZip).To(BeAnExistingFile())

			Expect(subject.Process(exampleZip, outputPath)).To(MatchError("Unable to find image in ZIP: alpha.png"))
		})
	})
})
