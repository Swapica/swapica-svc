package evm

import (
	"github.com/Swapica/swapica-svc/resources"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
)

func (e *evmProxy) GetMatchStatus(id *big.Int) (resources.Status, error) {
	matchStatus, err := e.swapper.MatchStatus(&bind.CallOpts{}, id)
	if err != nil {
		return resources.Status{}, err
	}

	result := resources.Status{
		State:      State(matchStatus.State),
		ExecutedBy: matchStatus.ExecutedBy,
	}

	return result, err
}