package fetcher


type Fetcher interface {
	Fetch(chan string)
}
