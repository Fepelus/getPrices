package fetcher

import (
	"github.com/Fepelus/getPrices/entities"
	"github.com/Fepelus/getPrices/outputter"
)

type Fetcher interface {
	Fetch(outputter.Outputter)
}

type fetcherConstructor func(string) Fetcher

var Brokers = map[string]fetcherConstructor{
	"YAHOO":    fetcherConstructor(NewYahoo),
	"VANGUARD": fetcherConstructor(NewVanguard),
	"QUANDL":   fetcherConstructor(NewQuandl),
	"BLOOMBERG":  fetcherConstructor(NewBloomberg),
	"EODDATA":  fetcherConstructor(NewEoddata),
	"FIXER":  fetcherConstructor(NewFixer),
}

func FetchAndOutput(commodities []entities.Commodity, output outputter.Outputter) {
	fetchers := makeFetchers(commodities)
	for _, thisfetcher := range fetchers {
		go thisfetcher.Fetch(output)
	}
}

func makeFetchers(commodities []entities.Commodity) []Fetcher {
	var output []Fetcher
	for _, commodity := range commodities {
		output = append(output, Brokers[commodity.Broker](commodity.Ticker))
	}
	return output
}
