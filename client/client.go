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
const imageFileParam = "image_file"
const bgImageFileParam = "bg_image_file"

//go:generate counterfeiter . ClientInterface
type ClientInterface interface {
	RemoveFromFile(inputPath string, apiKey string, params map[string]string) ([]byte, string, error)
}

type Client struct {
	Version    string
	HTTPClient http.Client
}

func (c Client) RemoveFromFile(inputPath string, apiKey string, params map[string]string) ([]byte, string, error) {
	request, err := c.buildRequest(APIEndpoint, apiKey, params, inputPath)
	if err != nil {
		return nil, "", err
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	contentType := resp.Header.Get("Content-Type")

	if statusCode == 200 {
		return body, contentType, err
	} else if statusCode >= 400 && statusCode < 500 {
		return nil, "", parseJsonErrors(statusCode, body)
	} else {
		return nil, "", fmt.Errorf("Unable to process image http_status=%d", statusCode)
	}
}

func (c Client) buildRequest(uri string, apiKey string, params map[string]string, inputPath string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err := attachFile(writer, imageFileParam, inputPath)
	if err != nil {
		return nil, err
	}

	if len(params[bgImageFileParam]) > 0 {
		err := attachFile(writer, bgImageFileParam, params[bgImageFileParam])
		if err != nil {
			return nil, err
		}
		delete(params, bgImageFileParam)
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
	req.Header.Add("User-Agent", c.userAgent())
	return req, err
}

func attachFile(writer *multipart.Writer, paramName string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("Unable to read file")
	}

	defer file.Close()

	part, err := writer.CreateFormFile(paramName, filepath.Base(filePath))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	return err
}

func (c Client) userAgent() string {
	return fmt.Sprintf("remove-bg-go-%s", c.Version)
}

func parseJsonErrors(statusCode int, body []byte) error {
	parsedErrorResponse := jsonErrorResponse{}
	err := json.Unmarshal(body, &parsedErrorResponse)
	if err != nil {
		return err
	}

	errorMessages := make([]string, len(parsedErrorResponse.Errors))
	for i, e := range parsedErrorResponse.Errors {
		errorMessages[i] = e.Title
	}

	message := strings.Join(errorMessages, ", ")

	return &RequestError{
		StatusCode: statusCode,
		Err:        errors.New(message),
	}
}

type jsonErrorResponse struct {
	Errors []struct {
		Title string
	}
}

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("%d: %s", r.StatusCode, r.Err.Error())
}

func (r *RequestError) RateLimitExceeded() bool {
	return r.StatusCode == 429
}
