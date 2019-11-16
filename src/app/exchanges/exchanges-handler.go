package exchanges

import "fmt"

type Exchanges interface {
	Start()
	GetHistoryTrades()
}

func Create(e *Exchanges) {
	fmt.Println(e)
}
