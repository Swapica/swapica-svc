package signature

import (
	sha3 "github.com/miguelmota/go-solidity-sha3"
)

type OrderData struct {
	OrderData string
}

func (log OrderData) Hash() []byte {
	return sha3.SoliditySHA3(log.OrderData)
}
