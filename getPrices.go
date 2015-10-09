package main

import (
	"github.com/Fepelus/getPrices/config"
	"github.com/Fepelus/getPrices/outputter"
	"github.com/Fepelus/getPrices/fetcher"
)

func main() {
	commodities := config.Parse()
	outputter := outputter.NewLedgerOutputter(len(commodities))
	fetcher.FetchAndOutput(commodities, outputter)
	outputter.Output()
}
