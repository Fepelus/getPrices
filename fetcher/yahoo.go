package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type yahoo string

func NewYahoo(label string) Fetcher {
	return yahoo(label)
}

func (this yahoo) Fetch(output chan string) {
	url := this.url()
	// http://finance.yahoo.com/d/quotes.csv?s=RKN.AX&f=snd1t1l1

	csv := this.call(url)
	// "RKN.AX","RECKON FPO","8/6/2015","12:19pm",2.03

	formatted := this.format(csv)
	// P 2015/08/06 12:19:00 RKN $2.03

	output <- formatted
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

func (this yahoo) format(csv string) string {
	split := strings.Split(csv, ",")
	date, _ := time.Parse("\"02/1/2006\"", split[2])
	dtfmt := date.Format("2006/02/01")
	clock, _ := time.Parse("\"15:04pm\"", split[3])
	clfmt := clock.Format("15:04")
	symsplit := strings.Split(split[0], "\"")
	qtsplit := strings.Split(symsplit[1], ".")
	price := strings.TrimSpace(split[4])
	commodity := qtsplit[0]
	return fmt.Sprintf("P %s %s:00 %s $%s", dtfmt, clfmt, commodity, price)
}
