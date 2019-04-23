package processor_test

import (
	"errors"
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
		subject        processor.Processor
	)

	BeforeEach(func() {
		fakeClient = &clientfakes.FakeClientInterface{}
		fakeFileWriter = &processorfakes.FakeFileWriterInterface{}

		subject = processor.Processor{
			APIKey:     "api-key",
			Client:     fakeClient,
			FileWriter: fakeFileWriter,
		}
	})

	It("coordinates the HTTP request and writing the result", func() {
		fakeClient.RemoveFromFileReturnsOnCall(0, []byte("Processed1"), nil)
		fakeClient.RemoveFromFileReturnsOnCall(1, []byte("Processed2"), nil)

		inputPaths := []string{"dir/image1.jpg", "dir/image2.jpg"}
		s := processor.Settings{
			OutputDirectory: "output-dir",
		}

		subject.Process(inputPaths, s)

		Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(2))

		clientArg1, clientArg2, params := fakeClient.RemoveFromFileArgsForCall(0)
		Expect(clientArg1).To(Equal("dir/image1.jpg"))
		Expect(clientArg2).To(Equal("api-key"))
		Expect(len(params)).To(Equal(0))

		Expect(fakeFileWriter.WriteCallCount()).To(Equal(2))

		writerArg1, writerArg2 := fakeFileWriter.WriteArgsForCall(0)
		Expect(writerArg1).To(Equal("output-dir/image1.png"))
		Expect(writerArg2).To(Equal([]byte("Processed1")))
	})

	It("passes non-empty image options to the client", func() {
		fakeClient.RemoveFromFileReturnsOnCall(0, []byte("Processed1"), nil)

		inputPaths := []string{"dir/image1.jpg"}
		s := processor.Settings{
			OutputDirectory: "output-dir",
			ImageSettings: processor.ImageSettings{
				Size:     "size-value",
				Type:     "type-value",
				Channels: "channels-value",
				BgColor:  "bg-color-value",
			},
		}

		subject.Process(inputPaths, s)

		Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(1))
		_, _, params := fakeClient.RemoveFromFileArgsForCall(0)

		Expect(params["size"]).To(Equal("size-value"))
		Expect(params["type"]).To(Equal("type-value"))
		Expect(params["channels"]).To(Equal("channels-value"))
		Expect(params["bg_color"]).To(Equal("bg-color-value"))
	})

	Context("client error", func() {
		It("keeps processing images", func() {
			fakeClient.RemoveFromFileReturnsOnCall(0, nil, errors.New("boom"))
			fakeClient.RemoveFromFileReturnsOnCall(1, []byte("Processed2"), nil)

			inputPaths := []string{"dir/image1.jpg", "dir/image2.jpg"}
			s := processor.Settings{
				OutputDirectory: "output-dir",
			}

			subject.Process(inputPaths, s)

			Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(2))
			Expect(fakeFileWriter.WriteCallCount()).To(Equal(1))

			_, writerArg2 := fakeFileWriter.WriteArgsForCall(0)
			Expect(writerArg2).To(Equal([]byte("Processed2")))
		})
	})

	Context("writer error", func() {
		It("keeps processing images", func() {
			fakeClient.RemoveFromFileReturnsOnCall(0, []byte("Processed1"), nil)
			fakeClient.RemoveFromFileReturnsOnCall(1, []byte("Processed2"), nil)
			fakeFileWriter.WriteReturnsOnCall(0, errors.New("boom"))

			inputPaths := []string{"dir/image1.jpg", "dir/image2.jpg"}
			s := processor.Settings{
				OutputDirectory: "output-dir",
			}

			subject.Process(inputPaths, s)

			Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(2))
		})
	})
})
