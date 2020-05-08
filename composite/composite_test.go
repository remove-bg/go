package composite_test

import (
	"crypto/sha256"
	"fmt"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/remove-bg/go/composite"
	"io"
	"os"
	"path"
	"runtime"
)

var _ = Describe("Composite", func() {
	var (
		subject       composite.Composite
		exampleZip    string
		referencePath string
		outputPath    string
		testDir       string
	)

	BeforeEach(func() {
		subject = composite.New()

		_, testFile, _, _ := runtime.Caller(0)
		testDir = path.Dir(testFile)

		exampleZip = path.Join(testDir, "../fixtures/zip/example-cat.zip")
		referencePath = path.Join(testDir, "../fixtures/zip/reference-example-cat.png")
		outputPath = path.Join(testDir, fmt.Sprintf("../tmp/composite-cat-%d.png", config.GinkgoConfig.ParallelNode))

		// Remove stale state from any previous test runs
		os.Remove(outputPath)
	})

	It("combines the color.jpg and alpha.png into a transparent PNG", func() {
		Expect(outputPath).ToNot(BeAnExistingFile())

		Expect(subject.Process(exampleZip, outputPath)).Should(Succeed())
		Expect(outputPath).To(BeAnExistingFile())

		outputSha := fileSha(outputPath)
		referenceSha := fileSha(referencePath)

		Expect(outputSha).To(Equal(referenceSha), "Expected output composite to match reference composite")
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

func fileSha(filepath string) []byte {
	Expect(filepath).To(BeAnExistingFile())

	f, err := os.Open(filepath)
	Expect(err).To(BeNil())

	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)

	Expect(err).To(BeNil())

	return h.Sum(nil)
}
