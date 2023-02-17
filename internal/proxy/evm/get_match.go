package evm

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (e *evmProxy) GetMatch(id *big.Int) (swapica.ISwapicaMatch, error) {
	offset := id.Sub(id, bigOne)
	matches, err := e.swapper.GetAllMatches(&bind.CallOpts{}, offset, bigOne)
	if err != nil {
		return swapica.ISwapicaMatch{}, err
	}
	if len(matches) == 0 {
		return swapica.ISwapicaMatch{}, errors.New("match order does not exist")
	}
	if len(matches) > 1 {
		panic(errNotSingleOrder)
	}

	return matches[0], nil
}
