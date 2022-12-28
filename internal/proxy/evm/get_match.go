package evm

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (e *evmProxy) GetMatch(id *big.Int) (swapica.SwapicaMatch, error) {
	match, err := e.swapper.Matches(&bind.CallOpts{}, id)
	if err != nil {
		return swapica.SwapicaMatch{}, err
	}

	result := swapica.SwapicaMatch{
		Id:            match.Id,
		OriginOrderId: match.OriginOrderId,
		Account:       match.Account,
		TokenToSell:   match.TokenToSell,
		AmountToSell:  match.AmountToSell,
		OriginChain:   match.OriginChain,
	}

	return result, err
}
