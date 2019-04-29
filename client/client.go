package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const APIEndpoint = "https://api.remove.bg/v1.0/removebg"
const Version = "1.0.0"

//go:generate counterfeiter . ClientInterface
type ClientInterface interface {
	RemoveFromFile(inputPath string, apiKey string, params map[string]string) ([]byte, error)
}

type Client struct {
	HTTPClient http.Client
}

func (c Client) RemoveFromFile(inputPath string, apiKey string, params map[string]string) ([]byte, error) {
	request, err := buildRequest(APIEndpoint, apiKey, params, inputPath)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unable to process image http_status=%d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func buildRequest(uri string, apiKey string, params map[string]string, inputPath string) (*http.Request, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, errors.New("Unable to read file")
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image_file", filepath.Base(inputPath))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-Api-Key", apiKey)
	req.Header.Add("User-Agent", userAgent())
	return req, err
}

func userAgent() string {
	return fmt.Sprintf("remove-bg-go-%s", Version)
}
