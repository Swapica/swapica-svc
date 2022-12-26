package evm

import (
	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (e *evmProxy) GetOrderStatus(id *big.Int) (swapica.SwapicaStatus, error) {
	orderStatus, err := e.swapper.OrderStatus(&bind.CallOpts{}, id)
	if err != nil {
		return swapica.SwapicaStatus{}, err
	}

	result := swapica.SwapicaStatus{
		State:        orderStatus.State,
		ExecutedBy:   orderStatus.ExecutedBy,
		MatchSwapica: orderStatus.MatchSwapica,
	}

	return result, err
}
