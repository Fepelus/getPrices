package main

import (
	"fmt"
	"os"

	"github.com/Fepelus/getPrices/config"
)

func main() {
	commodities := config.Parse()
	fetchers := config.MakeFetchers(commodities)
	receiver := make(chan string)
	for _, fetcher := range fetchers {
		go fetcher.Fetch(receiver)
	}
	for _ = range fetchers {
		fmt.Fprintln(os.Stdout, <-receiver)
	}
}
