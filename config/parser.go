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
      fmt.Fprintf(os.Stderr, "Could find no config file at %s\n", filename)
      os.Exit(1)
   }
	bytes, filereaderr := ioutil.ReadFile(filename)
	if filereaderr != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %s\n", filename, filereaderr)
      os.Exit(1)
	}
	return string(bytes)
}

func parseconfig(input string) []entities.Commodity {
	var output []entities.Commodity

	lines := strings.Split(input, "\n")
	for line_number, line := range lines {
     // ignore comments and blank lines
     if len(strings.TrimSpace(line)) == 0 || line[0] == '#' {
       continue;
     }

	  tokens := strings.Split(line, " ")
	  if _, ok := brokers[tokens[0]]; ok {
         output = append(output, entities.NewCommodity(tokens[0], tokens[1]))
     } else {
       fmt.Fprintf(os.Stderr, "Parser error: No broker found called '%s' in line %d: %s\n", tokens[0], line_number, line)
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
