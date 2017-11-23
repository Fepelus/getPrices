package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Fepelus/getPrices/entities"
	"github.com/Fepelus/getPrices/outputter"
)

type yahoo string

func NewYahoo(label string) Fetcher {
	return yahoo(label)
}

func (this yahoo) Fetch(output outputter.Outputter) {
	url := this.url()
	// http://finance.yahoo.com/d/quotes.csv?s=RKN.AX&f=snd1t1l1

	csv := this.call(url)
	// "RKN.AX","RECKON FPO","8/6/2015","12:19pm",2.03

	price := this.makePrice(csv)

	output.Append(price)
}

func (this yahoo) url() string {
	return fmt.Sprintf("http://finance.yahoo.com/d/quotes.csv?s=%s.AX&f=snd1t1l1", this)
}

func (this yahoo) call(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
	}
	return string(body)
}

func (this yahoo) makePrice(csv string) entities.Price {
	split := strings.Split(csv, ",")
    fmt.Println(split)
	date, _ := time.Parse("\"1/2/2006\"", split[2])
	clock, _ := time.Parse("\"15:04pm\"", split[3])
	symsplit := strings.Split(split[0], "\"")
	qtsplit := strings.Split(symsplit[1], ".")
	price := strings.TrimSpace(split[4])
	commodity := qtsplit[0]
	return entities.NewPrice(date, clock, commodity, price)
}
