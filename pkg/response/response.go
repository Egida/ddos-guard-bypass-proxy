package response

import (
	"io"
	"net/http"

	"github.com/birros/ddos-guard-bypass-proxy/pkg/common"
)

type HTTPResponse struct {
	StatusCode int
	Status     string
	Header     map[string][]string
	Body       io.ReadCloser
}

type HTTPResponseDTO struct {
	StatusCode int                 `json:"status_code,omitempty"`
	Status     string              `json:"status,omitempty"`
	Header     map[string][]string `json:"header,omitempty"`
	Body       string              `json:"body,omitempty"`
}

func newHTTPResponseDTO(res HTTPResponse) (HTTPResponseDTO, error) {
	body, err := common.MarshalBody(res.Body)
	if err != nil {
		return HTTPResponseDTO{}, err
	}

	r := HTTPResponseDTO{
		StatusCode: res.StatusCode,
		Status:     res.Status,
		Header:     common.CopyHeader(res.Header),
		Body:       body,
	}

	return r, nil
}

func parseHTTPResponseDTO(res HTTPResponseDTO) (HTTPResponse, error) {
	body, err := common.UnmarshalBody(res.Body)
	if err != nil {
		return HTTPResponse{}, err
	}

	r := HTTPResponse{
		StatusCode: res.StatusCode,
		Status:     res.Status,
		Header:     common.CopyHeader(res.Header),
		Body:       body,
	}

	return r, nil
}

func NewHTTPResponseDTO(res *http.Response) (HTTPResponseDTO, error) {
	r := HTTPResponse{
		StatusCode: res.StatusCode,
		Status:     res.Status,
		Header:     common.CopyHeader(res.Header),
		Body:       res.Body,
	}

	return newHTTPResponseDTO(r)
}

func NewHTTPResponseFromDTO(dto HTTPResponseDTO) (*http.Response, error) {
	res, err := parseHTTPResponseDTO(dto)
	if err != nil {
		return nil, err
	}

	r := &http.Response{
		StatusCode: res.StatusCode,
		Status:     res.Status,
		Header:     res.Header,
		Body:       res.Body,
	}

	return r, err
}
