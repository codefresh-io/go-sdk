package codefresh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	logger "github.com/izumin5210/gentleman-logger"
	"gopkg.in/h2non/gentleman.v2"
)

type (
	Codefresh interface {
		requestAPI(*requestOptions) (*http.Response, error)
		ITokenAPI
		IPipelineAPI
		IRuntimeEnvironmentAPI
	}
)

func New(opt *ClientOptions) Codefresh {
	client := gentleman.New()
	client.BaseURL(opt.Host)
	if opt.Debug == true {
		client.Use(logger.New(os.Stdout))
	}

	return &codefresh{
		host:   opt.Host,
		token:  opt.Auth.Token,
		client: &http.Client{},
	}
}

func (c *codefresh) requestAPI(opt *requestOptions) (*http.Response, error) {
	var body []byte
	finalURL := fmt.Sprintf("%s%s", c.host, opt.path)
	if opt.qs != nil {
		finalURL += "?"
		for k, v := range opt.qs {
			finalURL += fmt.Sprintf("%s=%s", k, v)
		}
	}
	if opt.body != nil {
		body, _ = json.Marshal(opt.body)
	}
	request, err := http.NewRequest(opt.method, finalURL, bytes.NewBuffer(body))
	request.Header.Set("Authorization", c.token)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (c *codefresh) decodeResponseInto(resp *http.Response, target interface{}) error {
	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *codefresh) getBodyAsString(resp *http.Response) (string, error) {
	body, err := c.getBodyAsBytes(resp)
	return string(body), err
}

func (c *codefresh) getBodyAsBytes(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
