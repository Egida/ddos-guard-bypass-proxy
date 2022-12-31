package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/birros/ddos-guard-bypass-proxy/pkg/client"
)

const proxyUrl = "http://127.0.0.1:8192/"

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://ipinfo.io/ip", nil)
	if err != nil {
		log.Panicln(err)
	}

	client := client.Client{
		URL: proxyUrl,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(res.StatusCode)

	io.Copy(os.Stdout, res.Body)
}
