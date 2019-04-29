package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

	statusCode := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)

	if statusCode == 200 {
		return body, err
	} else if statusCode >= 400 && statusCode < 500 {
		return nil, parseJsonErrors(body)
	} else {
		return nil, fmt.Errorf("Unable to process image http_status=%d", resp.StatusCode)
	}
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

func parseJsonErrors(body []byte) error {
	parsedErrorResponse := jsonErrorResponse{}
	err := json.Unmarshal(body, &parsedErrorResponse)
	if err != nil {
		return err
	}

	errorMessages := make([]string, len(parsedErrorResponse.Errors))
	for i, e := range parsedErrorResponse.Errors {
		errorMessages[i] = e.Title
	}

	return errors.New(strings.Join(errorMessages, ", "))
}

type jsonErrorResponse struct {
	Errors []struct {
		Title string
	}
}
