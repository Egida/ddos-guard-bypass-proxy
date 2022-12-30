package worker

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/SkYNewZ/go-flaresolverr"
	"github.com/birros/ddos-guard-bypass-proxy/pkg/request"
	"github.com/birros/ddos-guard-bypass-proxy/pkg/response"
	"github.com/google/uuid"
	lru "github.com/hashicorp/golang-lru/v2"
)

type payload struct {
	Url        string
	RequestDTO request.HTTPRequestDTO
}

type Worker struct {
	cache   *lru.Cache[string, response.HTTPResponseDTO]
	pending sync.Map
	once    sync.Once
	ch      chan payload
	c       *config
}

func New(
	cache *lru.Cache[string, response.HTTPResponseDTO],
	options ...Option,
) *Worker {
	c := parseOptions(options)

	return &Worker{
		cache:   cache,
		pending: sync.Map{},
		ch:      make(chan payload),
		c:       c,
	}
}

func (w *Worker) Background() {
	w.once.Do(w.background)
}

func (w *Worker) background() {
	for payload := range w.ch {
		w.process(payload)
	}
}

func (w *Worker) process(payload payload) {
	log.Println("DOING", payload.Url)
	defer log.Println("DONE", payload.Url)

	defer func() {
		w.pending.Delete(payload.Url)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), w.c.queryTimeout)
	defer cancel()

	var res *flaresolverr.Response
	client := flaresolverr.New(w.c.flareSolverrUrl, w.c.queryTimeout, nil)
	if payload.RequestDTO.Method == http.MethodGet {
		var err error
		res, err = client.Get(ctx, payload.RequestDTO.URL, uuid.Nil)
		if err != nil {
			log.Println(payload.Url, err)
			return
		}
	} else if payload.RequestDTO.Method == http.MethodPost {
		var err error
		res, err = client.Post(ctx, payload.RequestDTO.URL, uuid.Nil, payload.RequestDTO.Body)
		if err != nil {
			log.Println(payload.Url, err)
			return
		}
	} else {
		log.Println(payload.Url, http.ErrNotSupported)
		return
	}

	header := map[string][]string{
		"Status":                 {res.Solution.Headers.Status},
		"Date":                   {res.Solution.Headers.Date},
		"Content-Type":           {res.Solution.Headers.ContentType},
		"Expires":                {res.Solution.Headers.Expires},
		"Cache-Control":          {res.Solution.Headers.CacheControl},
		"Pragma":                 {res.Solution.Headers.Pragma},
		"X-Frame-Options":        {res.Solution.Headers.XFrameOptions},
		"X-Content-Type-Options": {res.Solution.Headers.XContentTypeOptions},
		"Cf-Cache-Status":        {res.Solution.Headers.CfCacheStatus},
		"Expect-Ct":              {res.Solution.Headers.ExpectCt},
		"Report-To":              {res.Solution.Headers.ReportTo},
		"Nel":                    {res.Solution.Headers.Nel},
		"Server":                 {res.Solution.Headers.Server},
		"Cf-Ray":                 {res.Solution.Headers.CfRay},
		"Content-Encoding":       {res.Solution.Headers.ContentEncoding},
		"Alt-Svc":                {res.Solution.Headers.AltSvc},
	}

	body := io.NopCloser(strings.NewReader(res.Solution.Response))

	r := &http.Response{
		StatusCode: res.Solution.Status,
		Status:     res.Status,
		Header:     header,
		Body:       body,
	}

	resDTO, err := response.NewHTTPResponseDTO(r)
	if err != nil {
		log.Println(payload.Url, http.ErrNotSupported)
		return
	}

	w.cache.Add(payload.Url, resDTO)
}

func (w *Worker) Add(req *http.Request) {
	_, ok := w.pending.Load(req.URL.String())
	if ok {
		log.Println("ALREADY", req.URL)
		return
	}

	log.Println("ADDED", req.URL)

	w.pending.Store(req.URL.String(), true)

	reqDTO, err := request.NewHTTPRequestDTO(req)
	if err != nil {
		log.Println(req.URL)
		return
	}

	w.ch <- payload{
		Url:        reqDTO.URL,
		RequestDTO: reqDTO,
	}
}

func (w *Worker) Get(req *http.Request) (*http.Response, error) {
	resDTO, ok := w.cache.Get(req.URL.String())
	if !ok {
		return nil, errors.New("not cached")
	}

	return response.NewHTTPResponseFromDTO(resDTO)
}
