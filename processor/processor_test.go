package processor_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/remove-bg/go/client/clientfakes"
	"github.com/remove-bg/go/processor"
	"github.com/remove-bg/go/processor/processorfakes"
	"github.com/remove-bg/go/storage/storagefakes"
)

var _ = Describe("Processor", func() {
	var (
		fakeClient   *clientfakes.FakeClientInterface
		fakeStorage  *storagefakes.FakeStorageInterface
		fakePrompt   *processorfakes.FakePromptInterface
		fakeNotifier *processorfakes.FakeNotifierInterface
		subject      processor.Processor
		testSettings processor.Settings
	)

	BeforeEach(func() {
		fakeClient = &clientfakes.FakeClientInterface{}
		fakeStorage = &storagefakes.FakeStorageInterface{}
		fakePrompt = &processorfakes.FakePromptInterface{}
		fakeNotifier = &processorfakes.FakeNotifierInterface{}
		fakePrompt.ConfirmLargeBatchReturns(true)
		fakeStorage.ExpandPathsStub = func(input []string) ([]string, error) {
			return input, nil
		}

		subject = processor.Processor{
			APIKey:   "api-key",
			Client:   fakeClient,
			Storage:  fakeStorage,
			Prompt:   fakePrompt,
			Notifier: fakeNotifier,
		}

		testSettings = processor.Settings{
			OutputDirectory:            "output-dir",
			LargeBatchConfirmThreshold: 50,
			ReprocessExisting:          false,
		}
	})

	It("expands globs in the input paths", func() {
		fakeStorage.ExpandPathsStub = func(input []string) ([]string, error) {
			return []string{"dir/image1.jpg"}, nil
		}

		subject.Process([]string{"dir/*.jpg"}, testSettings)

		Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(1))
		Expect(fakeStorage.ExpandPathsCallCount()).To(Equal(1))

		clientArg1, _, _ := fakeClient.RemoveFromFileArgsForCall(0)
		Expect(clientArg1).To(Equal("dir/image1.jpg"))
	})

	It("coordinates the HTTP request and writing the result", func() {
		fakeClient.RemoveFromFileReturnsOnCall(0, []byte("Processed1"), nil)
		fakeClient.RemoveFromFileReturnsOnCall(1, []byte("Processed2"), nil)

		inputPaths := []string{"dir/image1.jpg", "dir/image2.jpg"}

		subject.Process(inputPaths, testSettings)

		Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(2))

		clientArg1, clientArg2, params := fakeClient.RemoveFromFileArgsForCall(0)
		Expect(clientArg1).To(Equal("dir/image1.jpg"))
		Expect(clientArg2).To(Equal("api-key"))
		Expect(len(params)).To(Equal(0))

		Expect(fakeStorage.WriteCallCount()).To(Equal(2))

		writerArg1, writerArg2 := fakeStorage.WriteArgsForCall(0)
		Expect(writerArg1).To(Equal("output-dir/image1.png"))
		Expect(writerArg2).To(Equal([]byte("Processed1")))
	})

	Describe("image options", func() {
		It("passes non-empty image options to the client", func() {
			fakeClient.RemoveFromFileReturnsOnCall(0, []byte("Processed1"), nil)
			inputPaths := []string{"dir/image1.jpg"}

			testSettings.ImageSettings = processor.ImageSettings{
				Size:        "size-value",
				Type:        "type-value",
				Channels:    "channels-value",
				BgColor:     "bg-color-value",
				BgImageFile: "bg-image-file-value",
				Format:      "format-value",
			}

			subject.Process(inputPaths, testSettings)

			Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(1))
			_, _, params := fakeClient.RemoveFromFileArgsForCall(0)

			Expect(params["size"]).To(Equal("size-value"))
			Expect(params["type"]).To(Equal("type-value"))
			Expect(params["channels"]).To(Equal("channels-value"))
			Expect(params["bg_color"]).To(Equal("bg-color-value"))
			Expect(params["bg_image_file"]).To(Equal("bg-image-file-value"))
			Expect(params["format"]).To(Equal("format-value"))
		})

		It("parses any extra API options into params", func() {
			fakeClient.RemoveFromFileReturnsOnCall(0, []byte("Processed1"), nil)
			inputPaths := []string{"dir/image1.jpg"}

			testSettings.ImageSettings = processor.ImageSettings{
				Size:            "size-value",
				ExtraApiOptions: "option1=val1&option2=val2",
			}

			subject.Process(inputPaths, testSettings)

			Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(1))
			_, _, params := fakeClient.RemoveFromFileArgsForCall(0)

			Expect(params["size"]).To(Equal("size-value"))
			Expect(params["option1"]).To(Equal("val1"))
			Expect(params["option2"]).To(Equal("val2"))
		})
	})

	Context("client error", func() {
		It("keeps processing images", func() {
			fakeClient.RemoveFromFileReturnsOnCall(0, nil, errors.New("boom"))
			fakeClient.RemoveFromFileReturnsOnCall(1, []byte("Processed2"), nil)
			inputPaths := []string{"dir/image1.jpg", "dir/image2.jpg"}

			subject.Process(inputPaths, testSettings)

			Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(2))
			Expect(fakeStorage.WriteCallCount()).To(Equal(1))
			Expect(fakeNotifier.ErrorCallCount()).To(Equal(1))
			Expect(fakeNotifier.SuccessCallCount()).To(Equal(1))

			_, writerArg2 := fakeStorage.WriteArgsForCall(0)
			Expect(writerArg2).To(Equal([]byte("Processed2")))
		})

		It("passes the error details to the notifier", func() {
			err := errors.New("boom")
			fakeClient.RemoveFromFileReturnsOnCall(0, nil, err)
			fakeClient.RemoveFromFileReturnsOnCall(1, []byte("Processed2"), nil)
			inputPaths := []string{"dir/image1.jpg", "dir/image2.jpg"}

			subject.Process(inputPaths, testSettings)

			Expect(fakeNotifier.ErrorCallCount()).To(Equal(1))

			notifiedErr, notifiedPath, notifiedImageNumber, notifiedTotal := fakeNotifier.ErrorArgsForCall(0)

			Expect(notifiedErr).To(Equal(err))
			Expect(notifiedPath).To(Equal("dir/image1.jpg"))
			Expect(notifiedImageNumber).To(Equal(1))
			Expect(notifiedTotal).To(Equal(2))
		})
	})

	Context("writer error", func() {
		It("keeps processing images", func() {
			fakeClient.RemoveFromFileReturnsOnCall(0, []byte("Processed1"), nil)
			fakeClient.RemoveFromFileReturnsOnCall(1, []byte("Processed2"), nil)
			fakeStorage.WriteReturnsOnCall(0, errors.New("boom"))
			inputPaths := []string{"dir/image1.jpg", "dir/image2.jpg"}

			subject.Process(inputPaths, testSettings)

			Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(2))
			Expect(fakeNotifier.ErrorCallCount()).To(Equal(1))
			Expect(fakeNotifier.SuccessCallCount()).To(Equal(1))
		})

		It("passes the error details to the notifier", func() {
			err := errors.New("boom")
			fakeClient.RemoveFromFileReturnsOnCall(0, []byte("Processed1"), nil)
			fakeClient.RemoveFromFileReturnsOnCall(1, []byte("Processed2"), nil)
			fakeStorage.WriteReturnsOnCall(0, err)
			inputPaths := []string{"dir/image1.jpg", "dir/image2.jpg"}

			subject.Process(inputPaths, testSettings)

			Expect(fakeNotifier.ErrorCallCount()).To(Equal(1))

			notifiedErr, notifiedPath, notifiedImageNumber, notifiedTotal := fakeNotifier.ErrorArgsForCall(0)

			Expect(notifiedErr).To(Equal(err))
			Expect(notifiedPath).To(Equal("dir/image1.jpg"))
			Expect(notifiedImageNumber).To(Equal(1))
			Expect(notifiedTotal).To(Equal(2))
		})
	})

	Describe("skipping already processed files", func() {
		It("skips processing if the output file exists", func() {
			testSettings.OutputDirectory = ""
			inputPaths := []string{"dir/image1.jpg", "dir/image2.jpg"}
			fakeStorage.FileExistsReturnsOnCall(0, true)
			fakeStorage.FileExistsReturnsOnCall(1, false)

			subject.Process(inputPaths, testSettings)

			Expect(fakeNotifier.SkipCallCount()).To(Equal(1))
			Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(1))

			notifiedInput, notifiedOutput, notifiedImageNumber, notifiedTotal := fakeNotifier.SkipArgsForCall(0)
			Expect(notifiedInput).To(Equal("dir/image1.jpg"))
			Expect(notifiedOutput).To(Equal("dir/image1-removebg.png"))
			Expect(notifiedImageNumber).To(Equal(1))
			Expect(notifiedTotal).To(Equal(2))

			processedPath, _, _ := fakeClient.RemoveFromFileArgsForCall(0)
			Expect(processedPath).To(Equal("dir/image2.jpg"))
		})

		It("can be configured to force re-processing of existing files", func() {
			testSettings.ReprocessExisting = true
			fakeStorage.FileExistsReturns(true)

			subject.Process([]string{"dir/image1.jpg"}, testSettings)

			Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(1))
			Expect(fakeStorage.WriteCallCount()).To(Equal(1))
			Expect(fakeNotifier.SuccessCallCount()).To(Equal(1))
		})
	})

	Describe("large batch confirmation", func() {
		It("doesn't prompt under the limit", func() {
			inputPaths := []string{"dir/image1.jpg"}
			testSettings.LargeBatchConfirmThreshold = 50

			subject.Process(inputPaths, testSettings)

			Expect(fakePrompt.ConfirmLargeBatchCallCount()).To(Equal(0))
		})

		It("delegates to the prompt", func() {
			fakePrompt.ConfirmLargeBatchReturns(true)
			inputPaths := make([]string, 50)
			testSettings.LargeBatchConfirmThreshold = 50

			subject.Process(inputPaths, testSettings)

			Expect(fakePrompt.ConfirmLargeBatchCallCount()).To(Equal(1))
			Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(50))
		})

		It("can be skipped with a negative value", func() {
			inputPaths := make([]string, 50)
			testSettings.LargeBatchConfirmThreshold = -1

			subject.Process(inputPaths, testSettings)

			Expect(fakePrompt.ConfirmLargeBatchCallCount()).To(Equal(0))
			Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(50))
		})

		It("can allows configuration of the threshold", func() {
			fakePrompt.ConfirmLargeBatchReturns(true)
			inputPaths := make([]string, 25)
			testSettings.LargeBatchConfirmThreshold = 25

			subject.Process(inputPaths, testSettings)

			Expect(fakePrompt.ConfirmLargeBatchCallCount()).To(Equal(1))
			Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(25))
		})

		It("doesn't process if the confirmation is rejected", func() {
			fakePrompt.ConfirmLargeBatchReturns(false)
			inputPaths := make([]string, 50)
			testSettings.LargeBatchConfirmThreshold = 50

			subject.Process(inputPaths, testSettings)

			Expect(fakePrompt.ConfirmLargeBatchCallCount()).To(Equal(1))
			Expect(fakeClient.RemoveFromFileCallCount()).To(Equal(0))
		})
	})

	Describe("NewProcessor", func() {
		It("builds a processor", func() {
			p := processor.NewProcessor("api-key")

			Expect(p.APIKey).To(Equal("api-key"))
			Expect(p.Client).ToNot(BeNil())
			Expect(p.Storage).ToNot(BeNil())
			Expect(p.Prompt).ToNot(BeNil())
			Expect(p.Notifier).ToNot(BeNil())
		})
	})
})
