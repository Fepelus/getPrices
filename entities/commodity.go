package entities

type Commodity struct {
	Broker string
	Ticker string
}

func NewCommodity(broker string, ticker string) Commodity {
	return Commodity{broker, ticker}
}

type Fetcher interface {
	Fetch(chan string)
}
