package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Fepelus/getPrices/entities"
	"github.com/Fepelus/getPrices/outputter"
)

type fixer string

func NewFixer(label string) Fetcher {
	return fixer(label)
}

func (this fixer) Fetch(output outputter.Outputter) {
	url := this.url()
    // https://api.fixer.io/latest?base=AUD&symbols=EUR

	json := this.call(url)

	data := this.parseJson(json)

	output.Append(this.makePrice(data))
}

func (this fixer) url() string {
    return fmt.Sprintf("https://api.fixer.io/latest?base=AUD&symbols=%s", this)
}

func (this fixer) call(url string) string {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not fetch the fixer page: ", err)
		os.Exit(1)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

// Thank you https://mholt.github.io/json-to-go/
type FixerResponse struct {
    Base  string `json:"base"`
    Date  string `json:"date"`
    Rates map[string]float64 `json:"rates"`
}

func (this fixer) parseJson(markup string) FixerResponse {
	jsonv := []byte(markup)
	var f FixerResponse
	_ = json.Unmarshal(jsonv, &f)
	return f
}

func (this fixer) makePrice(data FixerResponse) entities.Price {
	valuestring := strconv.FormatFloat(data.Rates[string(this)], 'f', -1, 64)
	/*
		split := strings.Split(data.Dataset.Data[0][0].(string), ",")
		date, _ := time.Parse("\"1/2/2006\"", split[2])
		clock, _ := time.Parse("\"15:04pm\"", split[3])
		symsplit := strings.Split(split[0], "\"")
		qtsplit := strings.Split(symsplit[1], ".")
		price := strings.TrimSpace(split[4])
		commodity := qtsplit[0]
		return entities.NewPrice(date, clock, commodity, price)
	*/
	date, _ := time.Parse("2006-01-02", data.Date)
	clock, _ := time.Parse("15:04", "17:00")
	return entities.NewPrice(date, clock, string(this), valuestring)
}
