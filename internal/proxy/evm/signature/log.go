package signature

import (
	sha3 "github.com/miguelmota/go-solidity-sha3"
)

type CancelMatch struct {
	OrderData string
}

func (log CancelMatch) Hash() []byte {
	return sha3.SoliditySHA3(log.OrderData)
}
