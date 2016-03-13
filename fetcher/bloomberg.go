package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
    "strconv"
    "strings"
    "time"

	"github.com/Fepelus/getPrices/entities"
	"github.com/Fepelus/getPrices/outputter"
	"github.com/Jeffail/gabs"
)

type bloomberg string

func NewBloomberg(label string) Fetcher {
	return bloomberg(label)
}

func (this bloomberg) Fetch(output outputter.Outputter) {
	markup := this.call(this.url())

	price := this.makePrice(markup)

	output.Append(price)
}

func (this bloomberg) url() string {
	return fmt.Sprintf("http://www.bloomberg.com/quote/%s:AU", this)
	//return fmt.Sprintf("http://www.tradingroom.com.au/apps/qt/quote.ac?code=%s", this)
}

func (this bloomberg) call(url string) string {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not fetch the Bloomberg page: ", err)
		os.Exit(1)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

// To parse the JSON into:
/*
type MeasureType struct {
	MeasureCode string
}
type Pricedata struct {
	AsOfDate    string
	Price       float64
	MeasureType MeasureType
}

type Profile struct {
	Price []Pricedata
}

type BloombergJson struct {
	Fund_price []Profile
}
*/

// The markup contains a JSON object with all the data for Bloomberg's page
// and it is all on one line, so here I split the markup up into lines,
// scan until find the one that has the variable in it, chop off the start
// of the line which leaves just the JSON object which I can then parse
// and extract the useful data from.
func (this bloomberg) makePrice(markup string) entities.Price {
	for _, line := range strings.Split(markup, "\n") {
		if strings.Index(line, "bootstrappedData: ") > -1 {
			jsonv := []byte(line[18 : len(line)])
			jsonParsed, err := gabs.ParseJSON(jsonv)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not parse the Bloomberg page: ", err)
				os.Exit(1)
			}
            children, ok := jsonParsed.Children()
            if ok != nil {
				fmt.Fprintln(os.Stderr, "Could not parse the Bloomberg page")
				os.Exit(1)
            }
            price := children[0].Path("basicQuote.price").Data().(float64)
            date := children[0].Path("basicQuote.priceDate").Data().(string)
            parsedDate, _ := time.Parse("1/2/2006", date)
            timed := children[0].Path("basicQuote.priceTime").Data().(string)
            parsedTime, _ := time.Parse("3:04 PM", timed)
            newPrice := entities.NewPrice(parsedDate, parsedTime, string(this),
                strconv.FormatFloat(price, 'f', -1, 64),
            )
            return newPrice
		}
	}
	return entities.Price{}
}
