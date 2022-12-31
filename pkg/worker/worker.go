package worker

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/birros/ddos-guard-bypass-proxy/pkg/flaresolverr"
	"github.com/birros/ddos-guard-bypass-proxy/pkg/request"
	"github.com/birros/ddos-guard-bypass-proxy/pkg/response"
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
	session *flaresolverr.Session
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
		session: flaresolverr.New(c.flareSolverrUrl),
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

	var res *http.Response
	if payload.RequestDTO.Method != http.MethodGet {
		log.Println(payload.Url, http.ErrNotSupported)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), w.c.queryTimeout)
	defer cancel()

	res, err := w.session.Get(ctx, payload.RequestDTO.URL)
	if err != nil {
		log.Println(payload.Url, err)
		return
	}

	resDTO, err := response.NewHTTPResponseDTO(res)
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
		return nil, errors.New("key not found")
	}

	return response.NewHTTPResponseFromDTO(resDTO)
}
