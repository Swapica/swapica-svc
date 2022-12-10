package evm

import (
	"context"
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

const (
	TokenTypeNative = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
)

func GetTransactionOpts(tokenAddress string, sender common.Address, amount *big.Int) *bind.TransactOpts {
	if tokenAddress == TokenTypeNative {
		return buildTransactOptsWithValue(sender, amount)
	}
	return buildTransactOpts(sender)
}

func NewProxy(rpc string, signer signature.Signer, swapperContract string, confirmations int) (types.Proxy, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	chainID, err := client.ChainID(context.TODO())
	if err != nil {
		return nil, err
	}

	b, err := swapica.NewSwapica(common.HexToAddress(swapperContract), client)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create swapper contract for address %s", swapperContract)
	}

	return &evmProxy{
		client:          client,
		signer:          signer,
		chainID:         chainID,
		swapperContract: common.HexToAddress(swapperContract),
		swapper:         b,
		confirmations:   confirmations,
	}, nil
}

type evmProxy struct {
	client          *ethclient.Client
	signer          signature.Signer
	chainID         *big.Int
	swapperContract common.Address
	swapper         *swapica.Swapica
	confirmations   int
}
