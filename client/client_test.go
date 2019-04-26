package client_test

import (
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
		fixtureFile string
		subject     client.Client
	)

	BeforeEach(func() {
		_, testFile, _, _ := runtime.Caller(0)
		fixtureFile = path.Join(path.Dir(testFile), "../fixtures/person-in-field.jpg")
		subject = client.Client{
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
			BodyString("data")

		result, err := subject.RemoveFromFile(fixtureFile, "api-key", map[string]string{})

		Expect(err).To(Not(HaveOccurred()))
		Expect(result).To(Equal([]byte("data")))
		Expect(gock.IsDone()).To(BeTrue())
	})

	It("includes the client version", func() {
		gock.New("https://api.remove.bg").
			Post("/v1.0/removebg").
			MatchHeader("User-Agent", `^remove-bg-go-\d+\.\d+\.\d+$`).
			Reply(200).
			BodyString("data")

		subject.RemoveFromFile(fixtureFile, "api-key", map[string]string{})

		Expect(gock.IsDone()).To(BeTrue())
	})

	Context("HTTP error", func() {
		It("returns a clear error", func() {
			gock.New("https://api.remove.bg").
				Post("/v1.0/removebg").
				MatchHeader("X-Api-Key", "^api-key$").
				Reply(400)

			result, err := subject.RemoveFromFile(fixtureFile, "api-key", map[string]string{})

			Expect(result).To(BeNil())
			Expect(err).To(MatchError("Unable to process image http_status=400"))
		})
	})

	Context("input file doesn't exist", func() {
		It("returns a clear error", func() {
			nonExistentFile := "/tmp/not-a-file"
			result, err := subject.RemoveFromFile(nonExistentFile, "api-key", map[string]string{})

			Expect(result).To(BeNil())
			Expect(err).To(MatchError("Unable to read file"))
		})
	})
})
