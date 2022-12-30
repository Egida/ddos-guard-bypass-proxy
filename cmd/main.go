package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/birros/ddos-guard-bypass-proxy/pkg/server"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "The addr of the application.")
	flag.Parse()

	log.SetFlags(log.Lshortfile | log.Ltime)

	handler := &server.Server{
		Do: func(req *http.Request) (*http.Response, error) {
			log.Println(req.URL)
			return http.DefaultClient.Do(req)
		},
	}

	log.Println("Starting proxy server on", *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
