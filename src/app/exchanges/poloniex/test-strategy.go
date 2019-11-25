package poloniex

import (
	"fmt"
	"log"

	"github.com/shopspring/decimal"
)

var (
	s1          string
	BTC         = decimal.NewFromFloat(0.5)
	firstPrice  decimal.Decimal
	middlePrice decimal.Decimal
	thirdPrice  decimal.Decimal
	firstList   = []string{"USDC_BTC", "USDT_BTC"}
	secondList  = []string{"USDC_ETH", "USDT_ETH"}
)

func (p *PoloniexWorker) StrategyStart() {

	//change := "T"
	go func() {
		fmt.Println("Start strategy")
		for {
			t := <-p.WsTickers.Subs["TICKER"]
			if t.(WSTicker).Symbol == "USDC_BTC" || t.(WSTicker).Symbol == "USDT_BTC" {
				firstPrice = t.(WSTicker).Last
				//fmt.Println("BTC: ", firstPrice)
			} else if t.(WSTicker).Symbol == "USDC_ETH" || t.(WSTicker).Symbol == "USDT_ETH" {
				middlePrice = t.(WSTicker).Last
				//fmt.Println("USDC: ", middlePrice)
			} else if t.(WSTicker).Symbol == "BTC_ETH" {
				thirdPrice = t.(WSTicker).Last
				fmt.Printf("\t%v\r ", t.(WSTicker).Last)
			} else {
				continue
			}

			// s1 = "USD" + "%s"
			// prices, err := p.GetPrice(fmt.Sprintf(s1+change+"_BTC"), fmt.Sprintf(s1+change+"_ETH", change), "BTC_ETH")
			// if err != nil {
			// 	log.Println(err)
			// 	time.Sleep(5 * time.Second)
			// }
			//fmt.Println(firstPrice, middlePrice, thirdPrice)
			if !firstPrice.IsZero() && !middlePrice.IsZero() && !thirdPrice.IsZero() {
				profit := p.GetProfit(firstPrice, middlePrice, thirdPrice, decimal.NewFromFloat(0))
				if profit.Cmp(decimal.NewFromFloat(25.)) > 0. {
					fmt.Println("Triangular profit:", profit)
					log.Fatal()
				}
			}

		}
	}()
}

// Triangular strategy profit: 3 prices input and Fee
// Return profit without Fee.
func (p *PoloniexWorker) GetProfit(first, middle, third, Fee decimal.Decimal) decimal.Decimal {
	var profit decimal.Decimal
	if first.Div(middle).Mul(third).Cmp(decimal.NewFromFloat(1.)) > 0 {
		val := first.Mul(BTC)
		s := val.Div(middle)
		f := s.Mul(third)
		profit = f.Sub(BTC)
	}

	//diff := (arbprice.Div(first))

	// diff := diff.Sub(decimal.NewFromFloat(1.))
	// fmt.Println(diff)
	// diff = diff.Mul(decimal.NewFromFloat(100.))

	return profit.Mul(first)
}
