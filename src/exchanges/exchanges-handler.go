package exchanges

import "fmt"

type Exchanges interface {
	Start()
}

func Create(e *Exchanges) {
	fmt.Println(e)
	e.
}
