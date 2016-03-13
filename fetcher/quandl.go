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

type quandl string

func NewQuandl(label string) Fetcher {
	return quandl(label)
}

func (this quandl) Fetch(output outputter.Outputter) {
	url := this.url()
	// http://finance.yahoo.com/d/quotes.csv?s=RKN.AX&f=snd1t1l1

	json := this.call(url)

	data := this.parseJson(json)

	output.Append(this.makePrice(data))
}

func (this quandl) url() string {
	return fmt.Sprintf("https://www.quandl.com/api/v3/datasets/WGC/GOLD_DAILY_AUD.json?start_date=2016-01-01")
}

func (this quandl) call(url string) string {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not fetch the Propertydata page: ", err)
		os.Exit(1)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

// Thank you https://mholt.github.io/json-to-go/
type GoldData struct {
	Dataset struct {
		Data [][]interface{} `json:"data"`
	} `json:"dataset"`
}

func (this quandl) parseJson(markup string) GoldData {
	jsonv := []byte(markup)
	var f GoldData
	_ = json.Unmarshal(jsonv, &f)
	return f
}

func (this quandl) makePrice(data GoldData) entities.Price {
	datestring := data.Dataset.Data[0][0].(string)
	valuefloat := data.Dataset.Data[0][1].(float64)
	valuestring := strconv.FormatFloat(valuefloat, 'f', -1, 64)
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
	date, _ := time.Parse("2006-01-02", datestring)
	clock, _ := time.Parse("15:04", "17:00")
	return entities.NewPrice(date, clock, "GOLD", valuestring)
}
