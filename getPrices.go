package main

import (
	"fmt"
	"os"

	"github.com/Fepelus/getPrices/vanguard"
	"github.com/Fepelus/getPrices/yahoo"
)

type Fetcher interface {
	Fetch(chan string)
}

func main() {
	fetchers := []Fetcher{
		vanguard.NewVanguard("VAN"),
		yahoo.NewYahoo("IOO"),
		yahoo.NewYahoo("RGB"),
		yahoo.NewYahoo("RKN"),
		yahoo.NewYahoo("VAP"),
	}
	receiver := make(chan string)
	for _, fetcher := range fetchers {
		go fetcher.Fetch(receiver)
	}
	for _ = range fetchers {
		fmt.Fprintln(os.Stdout, <-receiver)
	}
}
