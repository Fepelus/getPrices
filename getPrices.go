package main

import (
	"github.com/Fepelus/getPrices/config"
   "github.com/Fepelus/getPrices/outputter"
)

func main() {
	commodities := config.Parse()
	fetchers := config.MakeFetchers(commodities)
   outputter := outputter.NewLedgerOutputter(len(fetchers))
   for _, fetcher := range fetchers {
		go fetcher.Fetch(outputter)
	}
   outputter.Output()
}
