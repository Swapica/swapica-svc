package evm

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/enums"
	"github.com/Swapica/swapica-svc/resources"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (e *evmProxy) GetOrderStatus(id *big.Int) (resources.Status, error) {
	orderStatus, err := e.swapper.OrderStatus(&bind.CallOpts{}, id)
	if err != nil {
		return resources.Status{}, err
	}

	result := resources.Status{
		State:      enums.State(orderStatus.State),
		ExecutedBy: orderStatus.ExecutedBy,
	}

	return result, err
}
