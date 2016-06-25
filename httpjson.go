package httpjson

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
}

func NewHttpClient(maxIdle int, requestTimeout int) *HttpClient {
	client := &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: maxIdle},
		Timeout:   time.Duration(requestTimeout) * time.Second,
	}
	return &HttpClient{client: client}
}

func (c *HttpClient) PostJson(url string, request interface{}, response interface{}) error {

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.NewDecoder(bytes.NewReader(respb)).Decode(response)
}
