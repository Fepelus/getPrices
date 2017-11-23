package fetcher

import (
	"fmt"
    "time"

	"github.com/Fepelus/getPrices/entities"
	"github.com/Fepelus/getPrices/outputter"
    "github.com/namsral/microdata"
)

type bloomberg string

func NewBloomberg(label string) Fetcher {
	return bloomberg(label)
}

func (this bloomberg) Fetch(output outputter.Outputter) {
    var data *microdata.Microdata
    data, _ = microdata.ParseURL(this.url())

	price := this.makePrice(data)

	output.Append(price)
}

func (this bloomberg) url() string {
	return fmt.Sprintf("http://www.bloomberg.com/quote/%s:AU", this)
	//return fmt.Sprintf("http://www.tradingroom.com.au/apps/qt/quote.ac?code=%s", this)
}


func (this bloomberg) makePrice(data *microdata.Microdata) entities.Price {
    for i := 0; i < len(data.Items) ; i++ {
        if (data.Items[i].Types[0] == "http://schema.org/Intangible/FinancialQuote") {
            price := data.Items[i].Properties["price"][0].(string)
            quoteTime := data.Items[i].Properties["quoteTime"][0].(string)
            parsedTime, _ := time.Parse(time.RFC3339, quoteTime)
            return entities.NewPrice(parsedTime.Local(), parsedTime.Local(), string(this), price)
        }
    }
	return entities.Price{}
}
