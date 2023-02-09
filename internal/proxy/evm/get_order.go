package evm

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (e *evmProxy) GetOrder(id *big.Int) (swapica.ISwapicaOrder, error) {
	orders, err := e.swapper.GetAllOrders(&bind.CallOpts{}, id, bigOne)
	if err != nil {
		return swapica.ISwapicaOrder{}, err
	}
	if 1 != len(orders) {
		panic(errNotSingleOrder)
	}

	return orders[0], err
}
