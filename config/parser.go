package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Fepelus/getPrices/entities"
	"github.com/Fepelus/getPrices/fetcher"
)

type fetcherConstructor func(string) fetcher.Fetcher

var brokers = map[string]fetcherConstructor{
  "YAHOO":    fetcherConstructor(fetcher.NewYahoo),
  "VANGUARD": fetcherConstructor(fetcher.NewVanguard),
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
   if _, err := os.Stat(filename); os.IsNotExist(err) {
      fmt.Printf("Could find no config file at %s\n", filename)
      os.Exit(1)
   }
	bytes, filereaderr := ioutil.ReadFile(filename)
	if filereaderr != nil {
		fmt.Printf("Error reading %s: %s\n", filename, filereaderr)
      os.Exit(1)
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

func MakeFetchers(config []entities.Commodity) []fetcher.Fetcher {
	var output []fetcher.Fetcher
	for _, commodity := range config {
     output = append(output, brokers[commodity.Broker](commodity.Ticker))
	}
	return output
}
