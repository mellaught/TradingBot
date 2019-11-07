package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	//model "tradingBot/src/exchanges/binance/model"
	models "tradingBot/src/models"

	"github.com/adshao/go-binance"
)

type BinanceWorker struct {
	Cli *binance.Client
}

// Create Binance worker: ApiKey, ApiSecret from conf
func CreateBinanceWorker(conf *models.ApiData) *BinanceWorker {

	return &BinanceWorker{
		Cli: binance.NewClient(conf.ApiKey, conf.ApiSecret),
	}
}

func WsFetchDepth(symbols []string, out chan<- *binance.WsPartialDepthEvent) {
	var latestAsks []binance.Ask
	wsPartialDepthHandler := func(event *binance.WsPartialDepthEvent) {
		latestAsks = event.Asks
		if len(latestAsks) > 0 {
			fmt.Println("Event processed", event.Symbol)
			out <- event
		}
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	symbolMap := make(map[string]string)
	for i := 0; i < len(symbols); i++ {
		symbolMap[symbols[i]] = "5"
	}
	doneC, stopC, err := binance.WsCombinedPartialDepthServe(symbolMap, wsPartialDepthHandler, errHandler)
	if err != nil {
		fmt.Println(err)
	}

	// use stopC to exit
	go func() {
		time.Sleep(40 * time.Second)
		stopC <- struct{}{}
	}()

	// remove this if you do not want to be blocked here
	<-doneC
}

func (b *BinanceWorker) GetHistoryTrades(start, end int64, number int) {
	trades, err := b.Cli.NewAggTradesService().
		Symbol("BTCUSDT").StartTime(start).EndTime(end).
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
