package evm

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (e *evmProxy) CreateOrder(params types.CreateOrderParams) (interface{}, error) {
	sender := common.HexToAddress(params.Sender)

	tx, err := e.createOrder(params, sender)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}
	confirmed := true
	return encodeTx(tx, sender, e.chainID, params.SrcChain.ID, &confirmed)
}

func (e *evmProxy) createOrder(params types.CreateOrderParams, sender common.Address) (*ethTypes.Transaction, error) {
	tokenToSell := common.HexToAddress(params.TokenToSell)
	tokenToBuy := common.HexToAddress(params.TokenToBuy)
	destChainId, err := parseChainParams(params.DestChain.ChainParams)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse chain params")
	}

	tx, err := e.swapper.CreateOrder(
		GetTransactionOpts(params.TokenToSell, sender, big.NewInt(int64(params.AmountToSell.Float()))),
		swapica.ISwapicaCreateOrderRequest{
			TokenToSell:      tokenToSell,
			AmountToSell:     big.NewInt(int64(params.AmountToSell.Float())),
			TokenToBuy:       tokenToBuy,
			AmountToBuy:      big.NewInt(int64(params.AmountToBuy.Float())),
			DestinationChain: big.NewInt(destChainId.ChainId),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tx")
	}

	return tx, nil
}
