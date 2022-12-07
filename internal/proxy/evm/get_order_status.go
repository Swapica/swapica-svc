package evm

import (
	"github.com/Swapica/swapica-svc/resources"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
)

func (e *evmProxy) GetOrderStatus(id *big.Int) (resources.Status, error) {
	orderStatus, err := e.swapper.OrderStatus(&bind.CallOpts{}, id)
	if err != nil {
		return resources.Status{}, err
	}

	result := resources.Status{
		State:      State(orderStatus.State),
		ExecutedBy: orderStatus.ExecutedBy,
	}

	return result, err
}
