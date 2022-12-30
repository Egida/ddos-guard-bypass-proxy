package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/birros/ddos-guard-bypass-proxy/pkg/client"
)

const proxyUrlStr = "http://127.0.0.1:8080/"

var proxyUrl *url.URL

func init() {
	var err error
	proxyUrl, err = url.Parse(proxyUrlStr)
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://ifconfig.me/", nil)
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
