package evm

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (e *evmProxy) GetMatch(id *big.Int) (swapica.ISwapicaMatch, error) {
	matches, err := e.swapper.GetAllMatches(&bind.CallOpts{}, id, bigOne)
	if err != nil {
		return swapica.ISwapicaMatch{}, err
	}
	if len(matches) != 1 {
		panic(errNotSingleOrder)
	}
	return matches[0], nil
}
