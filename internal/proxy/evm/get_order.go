package evm

import (
	"github.com/Swapica/swapica-svc/resources"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
)

func (e *evmProxy) GetOrder(id *big.Int) (resources.Order, error) {
	order, err := e.swapper.Orders(&bind.CallOpts{}, id)
	if err != nil {
		return resources.Order{}, err
	}

	result := resources.Order{
		Account:      order.Account,
		AmountToBuy:  order.AmountToBuy,
		AmountToSell: order.AmountToSell,
		DestChain:    order.DestChain,
		Id:           order.Id,
		TokenToBuy:   order.TokenToBuy,
		TokenToSell:  order.TokenToSell,
	}

	return result, err
}
