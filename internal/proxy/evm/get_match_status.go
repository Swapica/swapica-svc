package evm

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/enums"
	"github.com/Swapica/swapica-svc/resources"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (e *evmProxy) GetMatchStatus(id *big.Int) (resources.Status, error) {
	matchStatus, err := e.swapper.MatchStatus(&bind.CallOpts{}, id)
	if err != nil {
		return resources.Status{}, err
	}

	result := resources.Status{
		State:      enums.State(matchStatus.State),
		ExecutedBy: matchStatus.ExecutedBy,
	}

	return result, err
}
