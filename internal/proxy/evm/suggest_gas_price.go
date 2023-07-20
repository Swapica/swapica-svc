package evm

import (
	"context"
	"math/big"
)

func (e *evmProxy) SuggestGasPrice() (*big.Int, error) {
	return e.client.SuggestGasPrice(context.Background())
}
