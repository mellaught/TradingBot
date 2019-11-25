package exchanges

import "fmt"

type tickerKey struct {
	Base  string
	Quote string
}

type Ticker struct {
	BuyPrice  float64
	SellPrice float64
}

type Tickers struct {
	Currencies      []string
	currencyToIndex map[string]int

	tickers map[tickerKey]Ticker
}

// Exchange -- interface for all exchanges
type Exchange interface {
	Start()
	GetTickersFromChan() *Tickers
	GetHistoryTrades()
}

// Create exchange and start it
func Create(e *Exchange) {
	fmt.Println(e)
}

func GetTickets(e Exchange) *Tickers {
	return e.GetTickersFromChan()
}

