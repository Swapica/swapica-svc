package evm

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (e *evmProxy) GetOrder(id *big.Int) (swapica.ISwapicaOrder, error) {
	offset := id.Sub(id, bigOne)
	orders, err := e.swapper.GetAllOrders(&bind.CallOpts{}, offset, bigOne)
	if err != nil {
		return swapica.ISwapicaOrder{}, err
	}
	if len(orders) == 0 {
		return swapica.ISwapicaOrder{}, errors.New("order does not exist")
	}
	if len(orders) > 1 {
		panic(errNotSingleOrder)
	}

	return orders[0], err
}
