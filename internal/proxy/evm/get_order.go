package evm

import (
	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (e *evmProxy) GetOrder(id *big.Int) (swapica.SwapicaOrder, error) {
	order, err := e.swapper.Orders(&bind.CallOpts{}, id)
	if err != nil {
		return swapica.SwapicaOrder{}, err
	}

	result := swapica.SwapicaOrder{
		Id:           order.Id,
		Account:      order.Account,
		TokenToSell:  order.TokenToSell,
		TokenToBuy:   order.TokenToBuy,
		AmountToSell: order.AmountToSell,
		AmountToBuy:  order.AmountToBuy,
		DestChain:    order.DestChain,
	}

	return result, err
}
