package main_test

import (
	"crypto/sha256"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
)

var _ = Describe("Remove.bg CLI: Composite command", func() {
	var (
		exampleZip    string
		referencePath string
		outputPath    string
		testDir       string
		cliPath       string
	)

	BeforeSuite(func() {
		var err error
		cliPath, err = gexec.Build("github.com/remove-bg/go")

		log.Printf("CLI: %s", cliPath)

		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	BeforeEach(func() {
		_, testFile, _, _ := runtime.Caller(0)
		testDir = path.Dir(testFile)

		exampleZip = path.Join(testDir, "fixtures/zip/example-cat.zip")
		Expect(exampleZip).To(BeAnExistingFile())

		referencePath = path.Join(testDir, "fixtures/zip/reference-example-cat.png")
		Expect(referencePath).To(BeAnExistingFile())

		outputPath = path.Join(testDir, "tmp/composite-cat.png")

		// Remove stale state from any previous test runs
		os.Remove(outputPath)
		Expect(outputPath).ToNot(BeAnExistingFile())
	})

	It("combines the color.jpg and alpha.png into a transparent PNG", func() {
		command := exec.Command(cliPath, "zip2png", exampleZip, outputPath)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		Eventually(session, 30).Should(gexec.Exit())

		Expect(session.ExitCode()).To(Equal(0))
		Expect(session.Err).To(gbytes.Say("Processed zip"))
		Expect(outputPath).To(BeAnExistingFile())

		outputSha := fileSha(outputPath)
		referenceSha := fileSha(referencePath)

		Expect(outputSha).To(Equal(referenceSha), "Expected output composite to match reference composite")
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
