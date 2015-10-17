package outputter

import (
	"fmt"
	"os"

	"github.com/Fepelus/getPrices/entities"
)

type ledgerOutputter struct {
	channel chan entities.Price
	count   int
}

func NewLedgerOutputter(count int) Outputter {
	return ledgerOutputter{
		channel: make(chan entities.Price),
		count:   count,
	}
}

func (this ledgerOutputter) Append(input entities.Price) {
	this.channel <- input
}

func (this ledgerOutputter) Output() {
	for i := 0; i < this.count; i++ {
		this.formatPriceForLedger(<-this.channel)
	}
}

func (this ledgerOutputter) formatPriceForLedger(input entities.Price) {
	dtfmt := input.Date.Format("2006/01/02")
	clfmt := input.Clock.Format("15:04")
	fmt.Fprintf(os.Stdout, "P %s %s:00 %s $%s\n", dtfmt, clfmt, input.Ticker, input.Price)
}
