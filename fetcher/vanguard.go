package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Fepelus/getPrices/entities"
	"github.com/Fepelus/getPrices/outputter"
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

func (this vanguard) url() string {
    return "https://api.vanguard.com/rs/gre/gra/1.3.0/datasets/auw-retail-listview-data.jsonp"
}

func (this vanguard) call(url string) string {
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
type Fund struct {
	NavPrice string
	AsOfDate string
	Identifier string
}

type VanguardJson struct {
	FundData map[string]Fund
}

// The markup contains a JSON object with all the data for Vanguard's page
// and it is all on one line, so here I split the markup up into lines,
// scan until find the one that has the variable in it, chop off the start
// of the line which leaves just the JSON object which I can then parse
// and extract the useful data from.
func (this vanguard) makePrice(markup string) entities.Price {
		if strings.Index(markup, "callback") > -1 {
			jsonv := []byte(markup[9 : len(markup)-1])
			var f VanguardJson
			err := json.Unmarshal(jsonv, &f)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not parse the Vanguard page: ", err)
				os.Exit(1)
			}
			for _, funddata := range f.FundData {
					if funddata.Identifier == string(this) {
						date, _ := time.Parse(
                            "02 Jan 2006",
							funddata.AsOfDate,
						)
						return entities.NewPrice(
							date,
							time.Date(2009, time.November, 10, 17, 0, 0, 0, time.UTC),
							"VAN",
                            funddata.NavPrice[1 : len(funddata.NavPrice)],
							//strconv.FormatFloat(funddata.NavPrice, 'f', -1, 64),
						)
					}
            }
        }
	return entities.Price{}
}
