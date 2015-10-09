package entities

import "time"

type Price struct {
	Date time.Time
   Clock time.Time
	Ticker string
   Price string    // stored here in dollars but without a dollar sign
}

func NewPrice(indate time.Time, inclock time.Time, ticker string, price string) Price {
	return Price{indate, inclock, ticker, price}
}

