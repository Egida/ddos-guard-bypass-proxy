package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/birros/ddos-guard-bypass-proxy/pkg/request"
	"github.com/birros/ddos-guard-bypass-proxy/pkg/response"
)

var _ http.Handler = (*Server)(nil)

type Server struct {
	Do func(req *http.Request) (*http.Response, error)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var reqDto request.HTTPRequestDTO
	err = json.Unmarshal(reqData, &reqDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := request.NewHTTPRequestFromDTOWithContext(r.Context(), reqDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := s.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resDto, err := response.NewHTTPResponseDTO(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resData, err := json.Marshal(resDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resData)
}
