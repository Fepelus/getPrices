package fetcher

import "github.com/Fepelus/getPrices/outputter"

type Fetcher interface {
	Fetch(outputter.Outputter)
}
