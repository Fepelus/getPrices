package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	return "https://www.vanguardinvestments.com.au/retail/jsp/investments/retail?portId=8129##prices-and-distributions-tab"
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

type VanguardJson struct {
	Fund_price []Profile
}

// The markup contains a JSON object with all the data for Vanguard's page
// and it is all on one line, so here I split the markup up into lines,
// scan until find the one that has the variable in it, chop off the start
// of the line which leaves just the JSON object which I can then parse
// and extract the useful data from.
func (this vanguard) makePrice(markup string) entities.Price {
	for _, line := range strings.Split(markup, "\n") {
		if strings.Index(line, "pricenav =") > -1 {
			jsonv := []byte(line[30 : len(line)-2])
			var f VanguardJson
			err := json.Unmarshal(jsonv, &f)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not parse the Vanguard page: ", err)
				os.Exit(1)
			}
			for _, profile := range f.Fund_price {
				for _, pricedata := range profile.Price {
					if pricedata.MeasureType.MeasureCode == "SELL" {
						date, _ := time.Parse(
							"2006-01-02T15:04:05-07:00",
							pricedata.AsOfDate,
						)
						return entities.NewPrice(
							date,
							time.Date(2009, time.November, 10, 17, 0, 0, 0, time.UTC),
							string(this),
							strconv.FormatFloat(pricedata.Price, 'f', -1, 64),
						)
					}
				}
			}
		}
	}
	return entities.Price{}
}
