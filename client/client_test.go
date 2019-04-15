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
	AfterEach(func() {
		gock.Off()
	})

	It("requests the background removal", func() {
		_, testFile, _, _ := runtime.Caller(0)
		fixtureFile := path.Join(path.Dir(testFile), "../fixtures/person-in-field.jpg")

		gock.New("https://api.remove.bg").
			Post("/v1.0/removebg").
			MatchHeader("X-Api-Key", "^api-key$").
			Reply(200).
			BodyString("data")

		c := client.Client{
			HTTPClient: http.Client{},
		}

		result, err := c.RemoveFromFile(fixtureFile, "api-key", map[string]string{})

		Expect(err).To(Not(HaveOccurred()))
		Expect(result).To(Equal([]byte("data")))
		Expect(gock.IsDone()).To(BeTrue())
	})
})
