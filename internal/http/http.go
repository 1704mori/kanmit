package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    interface{}
}

func MakeRequest[R any](options RequestOptions) (R, error) {
	var reqBody io.Reader
	var response R

	if options.Body != nil {
		switch v := options.Body.(type) {
		case map[string]string:
			formData := url.Values{}
			for key, value := range v {
				formData.Add(key, value)
			}
			reqBody = strings.NewReader(formData.Encode())
		case string:
			reqBody = strings.NewReader(v)
		default:
			bodyJSON, err := json.Marshal(v)
			if err != nil {
				return response, err
			}
			reqBody = bytes.NewReader(bodyJSON)
		}
	}

	req, err := http.NewRequest(options.Method, options.URL, reqBody)
	if err != nil {
		return response, err
	}

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, err
	}

	return response, nil
}
