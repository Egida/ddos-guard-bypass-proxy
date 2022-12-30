package request

import (
	"context"
	"io"
	"net/http"

	"github.com/birros/ddos-guard-bypass-proxy/pkg/common"
)

type HTTPRequest struct {
	Method string
	URL    string
	Header map[string][]string
	Body   io.Reader
}

type HTTPRequestDTO struct {
	Method string              `json:"method,omitempty"`
	URL    string              `json:"url,omitempty"`
	Header map[string][]string `json:"header,omitempty"`
	Body   string              `json:"body,omitempty"`
}

func newHTTPRequestDTO(req HTTPRequest) (HTTPRequestDTO, error) {
	body, err := common.MarshalBody(req.Body)
	if err != nil {
		return HTTPRequestDTO{}, err
	}

	r := HTTPRequestDTO{
		Method: req.Method,
		URL:    req.URL,
		Header: common.CopyHeader(req.Header),
		Body:   body,
	}

	return r, nil
}

func parseHTTPRequestDTO(req HTTPRequestDTO) (HTTPRequest, error) {
	body, err := common.UnmarshalBody(req.Body)
	if err != nil {
		return HTTPRequest{}, err
	}

	r := HTTPRequest{
		Method: req.Method,
		URL:    req.URL,
		Header: common.CopyHeader(req.Header),
		Body:   body,
	}

	return r, nil
}

func NewHTTPRequestDTO(req *http.Request) (HTTPRequestDTO, error) {
	r := HTTPRequest{
		Method: req.Method,
		URL:    req.URL.String(),
		Header: common.CopyHeader(req.Header),
		Body:   req.Body,
	}

	return newHTTPRequestDTO(r)
}

func NewHTTPRequestFromDTOWithContext(
	ctx context.Context,
	dto HTTPRequestDTO,
) (*http.Request, error) {
	req, err := parseHTTPRequestDTO(dto)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, req.Method, req.URL, req.Body)
	if err != nil {
		return nil, err
	}

	r.Header = req.Header

	return r, err
}
