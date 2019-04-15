package processor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/remove-bg/go/client/clientfakes"
	"github.com/remove-bg/go/processor"
	"github.com/remove-bg/go/processor/processorfakes"
)

var _ = Describe("Processor", func() {
	var (
		fakeClient     *clientfakes.FakeClientInterface
		fakeFileWriter *processorfakes.FakeFileWriterInterface
	)

	BeforeEach(func() {
		fakeClient = &clientfakes.FakeClientInterface{}
		fakeFileWriter = &processorfakes.FakeFileWriterInterface{}
	})

	It("coordinates the HTTP request and writing the result", func() {
		fakeClient.RemoveFromFileReturnsOnCall(0, []byte("Processed1"), nil)
		fakeClient.RemoveFromFileReturnsOnCall(1, []byte("Processed2"), nil)

		subject := processor.Processor{
			APIKey:     "api-key",
			Client:     fakeClient,
			FileWriter: fakeFileWriter,
		}

		inputPaths := []string{"dir/image1.jpg", "dir/image2.jpg"}

		subject.Process(inputPaths, "output-dir")

		Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(2))

		clientArg1, clientArg2, _ := fakeClient.RemoveFromFileArgsForCall(0)
		Expect(clientArg1).To(Equal("dir/image1.jpg"))
		Expect(clientArg2).To(Equal("api-key"))

		Expect(fakeFileWriter.WriteCallCount()).To(Equal(2))

		writerArg1, writerArg2 := fakeFileWriter.WriteArgsForCall(0)
		Expect(writerArg1).To(Equal("output-dir/image1.png"))
		Expect(writerArg2).To(Equal([]byte("Processed1")))
	})
})
