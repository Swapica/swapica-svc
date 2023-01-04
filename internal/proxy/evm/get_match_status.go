package evm

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (e *evmProxy) GetMatchStatus(id *big.Int) (swapica.SwapicaStatus, error) {
	matchStatus, err := e.swapper.MatchStatus(&bind.CallOpts{}, id)
	if err != nil {
		return swapica.SwapicaStatus{}, err
	}

	result := swapica.SwapicaStatus{
		State:        matchStatus.State,
		ExecutedBy:   matchStatus.ExecutedBy,
		MatchSwapica: matchStatus.MatchSwapica,
	}

	return result, err
}
