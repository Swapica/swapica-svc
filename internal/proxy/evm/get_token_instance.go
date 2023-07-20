package evm

import (
	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/erc20"
	"github.com/ethereum/go-ethereum/common"
)

func (e *evmProxy) GetTokenInstance(address common.Address) (*erc20.Erc20, error) {
	return erc20.NewErc20(address, e.client)
}
