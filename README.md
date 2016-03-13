getPrices
=========

Tool for getting stock prices of my current portfolio.
It produces output appropriate for the
[ledger-cli](http://www.ledger-cli.org/)
tool that I use to track my finances. Certainly this
tool is for my own use.

Symbols to fetch and places to fetch them are set in a
config file.

Config file
---------

Put the name of your config file in the `GETPRICES_FILE`
environment variable. If this is not set then the tool
looks for `${HOME}/.portfoliorc` and if it finds nothing
there then it errors out.

The format of the file is as follows.

Empty lines or those with `'#'` at the first character
are ignored.

All other lines must be commodity lines. One commodity
line has the format:

`[YAHOO|VANGUARD|QUANDL|BLOOMBERG] [^ ]+`

That is, the start of the line must match the keys of a map
given in the file `/fetcher/fetcher.go` which currently is
`YAHOO`, `VANGUARD`, `QUANDL` or `BLOOMBERG`.  This is
followed by a single space character and a single ticker
identifier which is anything that is not a space character.
If there is a space after the ticker identifier then the
space is ignored along with everything on the line after it.

The first token is referred to in the code as a `broker`
which is not a good name: a better name would be something
like `PriceProvider` because this token tells the `fetcher`
code which method to use to get the price of the given
ticker identifier. `YAHOO` means to use the API at
`finance.yahoo.com` to get prices of commodities that are
traded on the ASX. `VANGUARD` means to scrape the
`vanguardinvestments.com.au` site to get the price of
their index fund (which is not traded on the ASX).
`QUANDL` means to use the commodities API provided by Quandl.
`BLOOMBERG` means to scrape the Bloomberg site.


Design
-----

![Message Sequencing Chart](https://raw.githubusercontent.com/Fepelus/getPrices/master/msc.png)
