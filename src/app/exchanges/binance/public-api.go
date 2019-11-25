package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/adshao/go-binance"
	"github.com/pkg/errors"
)

// Return Tickers and Exchange Info.
func (b *BinanceWorker) GetExchangeInfo() (error, []binance.Symbol) {
	info, err := b.Cli.NewExchangeInfoService().Do(context.TODO())
	if err != nil {
		return err, nil
	}

	return nil, info.Symbols
}

// Return Order Book
func (b *BinanceWorker) getOrderBook(symbol string, depth int) (book OrderBookInternal, err error) {
	orderBookURL, err := b.makeOrderBookURL(symbol, depth)
	if err != nil {
		return OrderBookInternal{}, errors.Wrapf(err, "could not make order book URL")
	}

	resp, err := http.Get(orderBookURL)
	if err != nil {
		return OrderBookInternal{}, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		time.Sleep(apiInterval)
	} else if resp.StatusCode != http.StatusOK {
		return OrderBookInternal{}, fmt.Errorf("getOrderBook received bad status code: %v", resp.StatusCode)
	}

	var data OrderBookResponse

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return OrderBookInternal{}, err
	}

	return SerializeBinanceOrderBookREST(data), nil
}

//
func (b *BinanceWorker) makeOrderBookURL(symbol string, depth int) (string, error) {
	u, err := url.Parse(depthURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("symbol", symbol)
	q.Set("limit", strconv.Itoa(depth))
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// Return History trades with currrent start and end time in UNIX
func (b *BinanceWorker) GetHistoryTrades(symbol string, start, end int64, number int) {
	trades, err := b.Cli.NewAggTradesService().
		Symbol(symbol).StartTime(start).EndTime(end).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	var sum int = 0
	var qul float64 = 0.
	var allQul float64 = 0.
	for _, t := range trades {
		i, _ := strconv.ParseFloat(t.Quantity, 64)
		allQul += i
		if t.IsBuyerMaker {
			sum++
			qul += i
		}

	}

	fmt.Println(float64(float64(sum)/float64(len(trades)))*100, float64(float64(qul)/float64(allQul))*100)
	// err = json.Unmarshal(jsonBlob, &rankings)
	// if err != nil {
	// 	// nozzle.printError("opening config file", err.Error())
	// }
	orders, err := json.Marshal(trades)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(orders))
	// rankingsJson, _ := json.Marshal(rankings)
	err = ioutil.WriteFile(fmt.Sprintf("%d.json", number), orders, 0644)
	// fmt.Printf("%+v", rankings)
}
