package poloniex

import (
	"encoding/json"
	"errors"
	"fmt"
)

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Error(msg string, args ...interface{}) error {
	if len(args) > 0 {
		return errors.New(fmt.Sprintf(msg, args))
	} else {
		return errors.New(msg)
	}
}

func (s subscription) toJSON() ([]byte, bool) {
	json, err := json.Marshal(s)
	if err != nil {
		return json, false
	}
	return json, true
}
