package storage_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	. "."
)

var _ = Describe("FileStorage", func() {
	var (
		subject FileStorage
		testDir string
	)

	BeforeEach(func() {
		subject = FileStorage{}

		_, testFile, _, _ := runtime.Caller(0)
		testDir = path.Dir(testFile)
	})

	Describe("FileExists", func() {
		It("is true when the file is present", func() {
			fixtureFile := path.Join(testDir, "../fixtures/person-in-field.jpg")

			Expect(subject.FileExists(fixtureFile)).To(BeTrue())
		})

		It("is false when the file doesn't exist", func() {
			missing := path.Join(testDir, "../fixtures/missing.jpg")

			Expect(subject.FileExists(missing)).To(BeFalse())
		})

		It("is false when the directory doesn't exist", func() {
			missingDir := path.Join(testDir, "../missing")

			Expect(subject.FileExists(missingDir)).To(BeFalse())
		})
	})

	Describe("ExpandPaths", func() {
		It("expands any star (*) globs in the inputs paths", func() {
			glob := path.Join(testDir, "../fixtures/*.jpg")
			expanded, err := subject.ExpandPaths([]string{glob})

			Expect(err).ToNot(HaveOccurred())
			Expect(expanded).To(ContainElement(MatchRegexp(`fixtures\/person-in-field\.jpg$`)))
		})

		It("expands any double-star (**) globs in the input paths", func() {
			glob := path.Join(testDir, "../fixtures/**/*.png")
			expanded, err := subject.ExpandPaths([]string{glob})

			Expect(err).ToNot(HaveOccurred())
			Expect(expanded).To(ContainElement(MatchRegexp(`nested\/plant\.png$`)))
		})

		It("expands any alternative patterns in the input paths", func() {
			glob := path.Join(testDir, "../fixtures/**/*.{jpg,png}")
			expanded, err := subject.ExpandPaths([]string{glob})

			Expect(err).ToNot(HaveOccurred())
			Expect(expanded).To(ContainElement(MatchRegexp(`fixtures\/person-in-field\.jpg$`)))
			Expect(expanded).To(ContainElement(MatchRegexp(`nested\/plant\.png$`)))
			Expect(expanded).ToNot(ContainElement(MatchRegexp(`nomatch\.txt$`)))
		})

		It("returns non-glob paths as-is", func() {
			fixtureFile := path.Join(testDir, "../fixtures/person-in-field.jpg")
			originals := []string{fixtureFile}
			expanded, err := subject.ExpandPaths(originals)

			Expect(err).ToNot(HaveOccurred())
			Expect(expanded).To(Equal(originals))
		})

		Context("input path isn't a glob", func() {
			// We want non-existent paths to remain, so we don't fail silently
			It("doesn't strip non-existent files", func() {
				inputPath := "missing/foo/bar.jpg"
				originals := []string{inputPath}
				expanded, err := subject.ExpandPaths(originals)

				Expect(err).ToNot(HaveOccurred())
				Expect(expanded).To(Equal(originals))
			})
		})
	})

	Describe("MkdirP", func() {
		var tmpDir string

		BeforeEach(func() {
			dir, err := ioutil.TempDir("", "mkdirp-spec")
			Expect(err).ToNot(HaveOccurred())

			tmpDir = dir
		})

		AfterEach(func() {
			os.RemoveAll(tmpDir)
		})

		It("creates deeply nested directories, if they don't exist", func() {
			outputDir := path.Join(tmpDir, "nested1/nested2")

			Expect(outputDir).ToNot(BeADirectory())
			Expect(subject.MkdirP(outputDir)).Should(Succeed())
			Expect(outputDir).To(BeADirectory())
		})
	})
})
