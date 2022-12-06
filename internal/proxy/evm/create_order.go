package evm

import (
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"math/big"
)

func (e *evmProxy) CreateOrder(params types.CreateOrderParams) (interface{}, error) {
	tokenType := TokenTypeErc20

	sender := common.HexToAddress(params.Sender)

	var tx *ethTypes.Transaction
	var err error

	switch tokenType {
	case TokenTypeErc20:
		tx, err = e.createOrderErc20(params, sender)
	default:
		return nil, errors.New("unsupported token type")
	}

	if err != nil {
		return nil, err
	}
	if tx == nil {
		// Token is already approved
		return nil, nil
	}

	return encodeTx(tx, sender, e.chainID, params.SrcChain.ID, nil)
}

func (e *evmProxy) createOrderErc20(params types.CreateOrderParams, sender common.Address) (*ethTypes.Transaction, error) {
	tokenToSell := common.HexToAddress(params.TokenToSell)
	tokenToBuy := common.HexToAddress(params.TokenToBuy)
	destChainId, err := parseChainParams(params.DestChain.ChainParams)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse chain params")
	}

	tx, err := e.swapper.CreateOrder(
		buildTransactOpts(sender),
		tokenToSell,
		params.AmountToSell.Int(),
		tokenToBuy,
		params.AmountToBuy.Int(),
		big.NewInt(destChainId.ChainId),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tx")
	}

	return tx, nil
}