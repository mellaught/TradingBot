package binance

import "github.com/adshao/go-binance"

func (b *BinanceWorker) makeErrorHandler() binance.ErrHandler {
	return func(err error) {
		b.log.Printf("Error in WS connection with Binance: %v", err)
	}
}
