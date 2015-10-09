package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Fepelus/getPrices/entities"
<<<<<<< HEAD
	"github.com/Fepelus/getPrices/fetcher"
)

type fetcherConstructor func(string) fetcher.Fetcher

var brokers = map[string]fetcherConstructor{
  "YAHOO":    fetcherConstructor(fetcher.NewYahoo),
  "VANGUARD": fetcherConstructor(fetcher.NewVanguard),
=======
	"github.com/Fepelus/getPrices/vanguard"
	"github.com/Fepelus/getPrices/yahoo"
)

var brokers = map[string]int{
	"YAHOO":    0,
	"VANGUARD": 1,
>>>>>>> df65aea0f670aced9332364d8485dec8dd44ec55
}

func Parse() []entities.Commodity {
	configstring := getConfigFromFile()
	return parseconfig(configstring)
}

func getConfigFromFile() string {
	filename := os.Getenv("GETPRICES_FILE")
	if filename == "" {
		home := os.Getenv("HOME")
		filename = fmt.Sprintf("%s/.portfoliorc", home)
	}
<<<<<<< HEAD
   if _, err := os.Stat(filename); os.IsNotExist(err) {
      fmt.Printf("Could find no config file at %s\n", filename)
      os.Exit(1)
   }
	bytes, filereaderr := ioutil.ReadFile(filename)
	if filereaderr != nil {
		fmt.Printf("Error reading %s: %s\n", filename, filereaderr)
      os.Exit(1)
=======
	bytes, filereaderr := ioutil.ReadFile(filename)
	if filereaderr != nil {
		fmt.Errorf("Error reading %s: %s\n", filename, filereaderr)
>>>>>>> df65aea0f670aced9332364d8485dec8dd44ec55
	}
	return string(bytes)
}

func parseconfig(input string) []entities.Commodity {
	var output []entities.Commodity

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if len(line) > 0 && line[0] != '#' {
			tokens := strings.Split(line, " ")
			if _, ok := brokers[tokens[0]]; ok {
				output = append(output, entities.NewCommodity(tokens[0], tokens[1]))
			} else {
				fmt.Errorf("Parser error: No broker found called %s", tokens[0])
			}
		}
	}

	return output
}

<<<<<<< HEAD
func MakeFetchers(config []entities.Commodity) []fetcher.Fetcher {
	var output []fetcher.Fetcher
	for _, commodity := range config {
     output = append(output, brokers[commodity.Broker](commodity.Ticker))
=======
func MakeFetchers(config []entities.Commodity) []entities.Fetcher {
	var output []entities.Fetcher
	for _, commodity := range config {
		if commodity.Broker == "YAHOO" {
			output = append(output, yahoo.NewYahoo(commodity.Ticker))
		}
		if commodity.Broker == "VANGUARD" {
			output = append(output, vanguard.NewVanguard(commodity.Ticker))
		}
>>>>>>> df65aea0f670aced9332364d8485dec8dd44ec55
	}
	return output
}
