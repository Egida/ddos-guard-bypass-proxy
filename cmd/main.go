package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/birros/ddos-guard-bypass-proxy/pkg/response"
	"github.com/birros/ddos-guard-bypass-proxy/pkg/server"
	"github.com/birros/ddos-guard-bypass-proxy/pkg/worker"
	lru "github.com/hashicorp/golang-lru/v2"
)

const (
	listenAddr = ":8192"
)

func main() {
	flag.Parse()

	log.SetFlags(log.Lshortfile | log.Ltime)

	cache, err := lru.New[string, response.HTTPResponseDTO](100)
	if err != nil {
		log.Panicln(err)
	}

	worker := worker.New(cache)
	go worker.Background()

	handler := &server.Server{
		Do: func(req *http.Request) (*http.Response, error) {
			log.Println(req.Method, req.URL)

			go worker.Add(req)
			return worker.Get(req)
		},
	}

	log.Println("starting proxy server on", listenAddr)
	err = http.ListenAndServe(listenAddr, handler)
	if err != nil {
		log.Panicln(err)
	}
}
