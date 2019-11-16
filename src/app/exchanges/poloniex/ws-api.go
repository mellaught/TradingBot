package poloniex

import (
	"fmt"

	polo "github.com/iowar/poloniex"
)

func (p *PoloniexWorker) SubscribeTikers() {
	err := p.WsCli.SubscribeTicker()
	if err != nil {
		return
	}
	for {
		t := <-p.WsCli.Subs["TICKER"]
		fmt.Println("SymBol: ", t.(polo.WSTicker).Symbol)
	}
}

func (p *PoloniexWorker) SubscribeMarkets(symbol string) {
	err := p.WsCli.SubscribeMarket(symbol)
	if err != nil {
		return
	}
	for {
		t := <-p.WsCli.Subs[symbol]
		fmt.Println("Market: ", t.([]polo.MarketUpdate)[0].Data)
	}
}



func (p *PoloniexWorker) SubscriPrices() {
	for _, symbol := range p.symbols {
		err := p.WsCli.SubscribePrice(symbol)
		if err != nil {
			return
		}
		for {
			t := <-p.WsCli.Subs[symbol]
			fmt.Println("Market: ", t.([]polo.MarketUpdate)[0].Data)
		}
	}
}
