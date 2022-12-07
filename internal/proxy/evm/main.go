package evm

import (
	"context"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"math/big"
)

const (
	TokenTypeNative = "native"
	TokenTypeErc20  = "erc20"
)

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

func (e *evmProxy) CancelOrder() (interface{}, error) {
	//TODO implement me
	panic("implement me")
}


func (e *evmProxy) ExecuteMatch() (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (e *evmProxy) ExecuteOrder() (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (e *evmProxy) GetMatch(matchId *big.Int) (struct {
	Id            *big.Int
	OriginOrderId *big.Int
	Account       common.Address
	TokenToSell   common.Address
	AmountToSell  *big.Int
	OriginChain   *big.Int
}, error) {
	match, err := e.swapper.Matches(&bind.CallOpts{}, matchId)
	if err != nil {
		return struct {
			Id            *big.Int
			OriginOrderId *big.Int
			Account       common.Address
			TokenToSell   common.Address
			AmountToSell  *big.Int
			OriginChain   *big.Int
		}{}, errors.Wrap(err, "failed to get match")
	}
	return match, nil
}

func (e *evmProxy) GetOrder(orderId *big.Int) (struct {
	Id           *big.Int
	Account      common.Address
	TokenToSell  common.Address
	TokenToBuy   common.Address
	AmountToSell *big.Int
	AmountToBuy  *big.Int
	DestChain    *big.Int
}, error) {
	order, err := e.swapper.Orders(&bind.CallOpts{}, orderId)
	if err != nil {
		return struct {
			Id           *big.Int
			Account      common.Address
			TokenToSell  common.Address
			TokenToBuy   common.Address
			AmountToSell *big.Int
			AmountToBuy  *big.Int
			DestChain    *big.Int
		}{}, errors.Wrap(err, "failed to get order")
	}
	return order, nil
}
