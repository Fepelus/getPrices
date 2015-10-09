package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
  
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
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

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

func (this vanguard) makePrice(markup string) entities.Price {
  for _, line := range strings.Split(markup, "\n") {
		if strings.Index(line, "jsonv =") > -1 {
			jsonv := []byte(line[12:])
			var f Quote
			_ = json.Unmarshal(jsonv, &f)
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

func (this vanguard) format(markup string) string {
	for _, line := range strings.Split(markup, "\n") {
		if strings.Index(line, "jsonv =") > -1 {
			jsonv := []byte(line[12:])
			var f Quote
			_ = json.Unmarshal(jsonv, &f)
			for _, fund := range f.Managedfund {
				if fund.Benchmark == "S&P/ASX 300 Index" {
					date, _ := time.Parse(
						"01/02/2006",
						fund.Unitpricedata.Effectivelatestdatetime,
					)
					dtfmt := date.Format("2006/02/01")
					return fmt.Sprintf("P %s 17:00:00 %s $%s", dtfmt, this, fund.Unitpricedata.Pricevaluelatestsell)
				}
			}
		}
	}
	return ""
}
