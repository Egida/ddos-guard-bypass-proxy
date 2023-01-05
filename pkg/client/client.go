package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/birros/ddos-guard-bypass-proxy/pkg/request"
	"github.com/birros/ddos-guard-bypass-proxy/pkg/response"
	"github.com/birros/ddos-guard-bypass-proxy/pkg/worker"
)

type Client struct {
	c *config
}

func NewClient(options ...Option) *Client {
	c := parseOptions(options)

	return &Client{
		c: c,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.c.useSession {
		req.Header.Add(worker.UseSessionHeader, "true")
	}

	reqDto, err := request.NewHTTPRequestDTO(req)
	if err != nil {
		return nil, err
	}

	reqData, err := json.Marshal(reqDto)
	if err != nil {
		return nil, err
	}

	reqProxy, err := http.NewRequestWithContext(
		req.Context(),
		"POST",
		c.c.proxyUrl,
		bytes.NewReader(reqData),
	)
	if err != nil {
		return nil, err
	}
	defer reqProxy.Body.Close()

	reqProxy.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	resProxy, err := client.Do(reqProxy)
	if err != nil {
		return nil, err
	}
	defer resProxy.Body.Close()

	resData, err := io.ReadAll(resProxy.Body)
	if err != nil {
		return nil, err
	}

	if resProxy.StatusCode != http.StatusOK {
		return nil, errors.New(string(resData))
	}

	var resDto response.HTTPResponseDTO
	err = json.Unmarshal(resData, &resDto)
	if err != nil {
		return nil, err
	}

	res, err := response.NewHTTPResponseFromDTO(resDto)
	if err != nil {
		return nil, err
	}

	return res, nil
}
