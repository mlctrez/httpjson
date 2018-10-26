package httpjson

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const ApplicationJson = "application/json"

type HttpJson struct {
	client *http.Client
}

func New(maxIdle int, requestTimeout int) *HttpJson {
	client := &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: maxIdle},
		Timeout:   time.Duration(requestTimeout) * time.Second,
	}
	return NewWithClient(client)
}

func NewWithClient(client *http.Client) *HttpJson {
	return &HttpJson{client: client}
}

func newRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if req != nil {
		req.Header.Set("Accept", ApplicationJson)
	}
	return
}

func (c *HttpJson) Post(url string, request interface{}, response interface{}) error {

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := newRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(response)
}

func (c *HttpJson) Get(url string, response interface{}) error {

	req, err := newRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(response)
}
