package client_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/remove-bg/go/client"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"path"
	"runtime"
)

var _ = Describe("Client", func() {
	var (
		fixtureFile   string
		bgFixtureFile string
		subject       client.Client
	)

	BeforeEach(func() {
		_, testFile, _, _ := runtime.Caller(0)
		fixtureFile = path.Join(path.Dir(testFile), "../fixtures/person-in-field.jpg")
		bgFixtureFile = path.Join(path.Dir(testFile), "../fixtures/background.jpg")
		subject = client.Client{
			Version:    "x.y.z",
			HTTPClient: http.Client{},
		}
	})

	AfterEach(func() {
		gock.Off()
	})

	It("requests the background removal", func() {
		gock.New("https://api.remove.bg").
			Post("/v1.0/removebg").
			MatchHeader("X-Api-Key", "^api-key$").
			Reply(200).
			SetHeader("Content-Type", "image/png").
			BodyString("data")

		result, contentType, err := subject.RemoveFromFile(fixtureFile, "api-key", map[string]string{})

		Expect(err).To(Not(HaveOccurred()))
		Expect(result).To(Equal([]byte("data")))
		Expect(contentType).To(Equal("image/png"))
		Expect(gock.IsDone()).To(BeTrue())
	})

	It("attaches the image file", func() {
		matcher := newMultipartAttachmentMatcher("image_file", "person-in-field.jpg")

		gock.New("https://api.remove.bg").
			Post("/v1.0/removebg").
			SetMatcher(matcher).
			Reply(200).
			BodyString("data")

		_, _, err := subject.RemoveFromFile(fixtureFile, "api-key", map[string]string{})

		Expect(err).To(Not(HaveOccurred()))
		Expect(gock.IsDone()).To(BeTrue())
	})

	It("attaches a background image file if specified", func() {
		imageMatcher := newMultipartAttachmentMatcher("image_file", "person-in-field.jpg")
		bgImagematcher := newMultipartAttachmentMatcher("bg_image_file", "background.jpg")

		gock.New("https://api.remove.bg").
			Post("/v1.0/removebg").
			SetMatcher(imageMatcher).
			SetMatcher(bgImagematcher).
			Reply(200).
			BodyString("data")

		params := map[string]string{
			"bg_image_file": bgFixtureFile,
		}

		_, _, err := subject.RemoveFromFile(fixtureFile, "api-key", params)

		Expect(err).To(Not(HaveOccurred()))
		Expect(gock.IsDone()).To(BeTrue())
	})

	It("includes the client version", func() {
		gock.New("https://api.remove.bg").
			Post("/v1.0/removebg").
			MatchHeader("User-Agent", "remove-bg-go-x.y.z").
			Reply(200).
			BodyString("data")

		subject.RemoveFromFile(fixtureFile, "api-key", map[string]string{})

		Expect(gock.IsDone()).To(BeTrue())
	})

	Context("server HTTP error", func() {
		It("returns a clear error", func() {
			gock.New("https://api.remove.bg").
				Post("/v1.0/removebg").
				Reply(500)

			result, _, err := subject.RemoveFromFile(fixtureFile, "api-key", map[string]string{})

			Expect(result).To(BeNil())
			Expect(err).To(MatchError("Unable to process image http_status=500"))
		})
	})

	Context("client HTTP error", func() {
		It("parses the JSON error messages", func() {
			jsonError := `{"errors": [{"title": "File too large"}, {"title": "Second error"}]}`

			gock.New("https://api.remove.bg").
				Post("/v1.0/removebg").
				Reply(400).
				BodyString(jsonError)

			result, _, err := subject.RemoveFromFile(fixtureFile, "api-key", map[string]string{})

			Expect(result).To(BeNil())
			Expect(err).To(MatchError("File too large, Second error"))
		})
	})

	Context("input file doesn't exist", func() {
		It("returns a clear error", func() {
			nonExistentFile := "/tmp/not-a-file"
			result, _, err := subject.RemoveFromFile(nonExistentFile, "api-key", map[string]string{})

			Expect(result).To(BeNil())
			Expect(err).To(MatchError("Unable to read file"))
		})
	})
})

func newMultipartAttachmentMatcher(key string, expectedFilename string) *gock.MockMatcher {
	// Create a new custom matcher with HTTP headers only matchers
	matcher := gock.NewBasicMatcher()

	// Add a custom match function
	matcher.Add(func(req *http.Request, ereq *gock.Request) (bool, error) {
		_, header, err := req.FormFile(key)
		if err != nil {
			return false, err
		}

		if header.Size == 0 {
			return false, fmt.Errorf("Attachment is empty: %v", header.Size)
		}

		if header.Filename == expectedFilename {
			return true, nil
		} else {
			return false, fmt.Errorf("Image filename was: %s", header.Filename)
		}
	})

	return matcher
}
