package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
    "time"

	"github.com/Fepelus/getPrices/entities"
	"github.com/Fepelus/getPrices/outputter"
)

type eoddata string

func NewEoddata(label string) Fetcher {
	return eoddata(label)
}

func (this eoddata) Fetch(output outputter.Outputter) {
	markup := this.call(this.url())

    //fmt.Printf(markup)

	price := this.makePrice(markup)

	output.Append(price)
}

func (this eoddata) url() string {
    return fmt.Sprintf("http://eoddata.com/stockquote/ASX/%s.htm", this)
}

func (this eoddata) call(url string) string {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not fetch the Eoddata page: ", err)
		os.Exit(1)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}


func (this eoddata) makePrice(markup string) entities.Price {
    splits  := strings.Split(markup, "Open Interest")
    lines := strings.Split(splits[1], "\n")
    tds := strings.Split(lines[1], "</")
    datestring :=  tds[0][strings.LastIndex(tds[0], ">")+1:]
    pricestring := tds[1][strings.LastIndex(tds[1], ">")+1:]
    date, _ := time.Parse( "01/02/06", datestring)
    fivepm := time.Date(2009, time.November, 10, 17, 0, 0, 0, time.UTC)
    return entities.NewPrice(
        date,
        fivepm,
        string(this),
        pricestring,
    )
}
