package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"os"

	"github.com/Fepelus/getPrices/outputter"
	"github.com/Fepelus/getPrices/entities"
)

type vanguard string

func NewVanguard(label string) Fetcher {
	return vanguard(label)
}

func (this vanguard) Fetch(output outputter.Outputter) {
	markup := this.call(this.url())

	price := this.makePrice(markup)

	output.Append(price)
}

func (this vanguard)url() string {
	return "https://www.vanguardinvestments.com.au/retail/ret/investments/managed-funds-retail.jsp"
}

func (this vanguard)call(url string) string {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not fetch the Vanguard page: ", err)
		os.Exit(1)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

// To parse the JSON into:
type Pricedata struct {
	Effectivelatestdatetime string
	Pricevaluelatestsell    string
}

type Fund struct {
	Benchmark     string
	Unitpricedata Pricedata
}

type Quote struct {
	Managedfund []Fund
}

// The markup contains a JSON object with all the data for Vanguard's page
// and it is all on one line, so here I split the markup up into lines,
// scan until find the one that has the variable in it, chop off the start
// of the line which leaves just the JSON object which I can then parse
// and extract the useful data from.
func (this vanguard) makePrice(markup string) entities.Price {
	for _, line := range strings.Split(markup, "\n") {
		if strings.Index(line, "jsonv =") > -1 {
			jsonv := []byte(line[12:])
			var f Quote
			err := json.Unmarshal(jsonv, &f)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not parse the Vanguard page: ", err)
				os.Exit(1)
			}
			for _, fund := range f.Managedfund {
				if fund.Benchmark == "S&P/ASX 300 Index" {
					date, _ := time.Parse(
						"01/02/2006",
						fund.Unitpricedata.Effectivelatestdatetime,
					)
					return entities.NewPrice(
						date,
						time.Date(2009, time.November, 10, 17, 0, 0, 0, time.UTC),
						string(this),
						fund.Unitpricedata.Pricevaluelatestsell,
					)
				}
			}
		}
	}
	return entities.Price{}
}
