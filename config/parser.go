package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Fepelus/getPrices/entities"
	"github.com/Fepelus/getPrices/vanguard"
	"github.com/Fepelus/getPrices/yahoo"
)

var brokers = map[string]int{
	"YAHOO":    0,
	"VANGUARD": 1,
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
	bytes, filereaderr := ioutil.ReadFile(filename)
	if filereaderr != nil {
		fmt.Errorf("Error reading %s: %s\n", filename, filereaderr)
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

func MakeFetchers(config []entities.Commodity) []entities.Fetcher {
	var output []entities.Fetcher
	for _, commodity := range config {
		if commodity.Broker == "YAHOO" {
			output = append(output, yahoo.NewYahoo(commodity.Ticker))
		}
		if commodity.Broker == "VANGUARD" {
			output = append(output, vanguard.NewVanguard(commodity.Ticker))
		}
	}
	return output
}
