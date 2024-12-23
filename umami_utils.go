package traefik_umami_feeder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func sendRequest(url string, body interface{}, headers http.Header) (*http.Response, error) {
	var req *http.Request
	var err error

	if body != nil {
		bodyJson, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(bodyJson))
	} else {
		req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	}

	if err != nil {
		return nil, err
	}

	if headers != nil {
		req.Header = headers
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	status := response.StatusCode
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("request failed with status %d", status)
	}

	return response, nil
}

func sendRequestAndParse(url string, body interface{}, headers http.Header, value interface{}) error {
	resp, err := sendRequest(url, body, headers)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, &value)
	if err != nil {
		return err
	}

	return nil
}
