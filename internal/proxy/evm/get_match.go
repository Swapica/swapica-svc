package evm

import (
	"github.com/Swapica/swapica-svc/resources"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
)

func (e *evmProxy) GetMatch(id *big.Int) (resources.Match, error) {
	match, err := e.swapper.Matches(&bind.CallOpts{}, id)
	if err != nil {
		return resources.Match{}, err
	}

	result := resources.Match{
		Account:       match.Account,
		AmountToSell:  match.AmountToSell,
		Id:            match.Id,
		OriginChain:   match.OriginChain,
		OriginOrderId: match.OriginOrderId,
		TokenToSell:   match.TokenToSell,
	}

	return result, err
}
