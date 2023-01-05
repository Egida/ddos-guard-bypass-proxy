package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/birros/ddos-guard-bypass-proxy/pkg/client"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://ipinfo.io/ip", nil)
	if err != nil {
		log.Panicln(err)
	}

	client := client.NewClient()
	res, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(res.StatusCode)

	io.Copy(os.Stdout, res.Body)
}
