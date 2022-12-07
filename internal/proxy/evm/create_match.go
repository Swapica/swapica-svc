package evm

import (
	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (e *evmProxy) CreateMatch(params types.CreateMatchParams) (interface{}, error) {
	srcChainParams, err := parseChainParams(params.SrcChain.ChainParams)
	if err != nil {
		return nil, err
	}

	calldata, err := CreateMatchCalldata(createMatchCalldata{
		Selector:     createMatch,
		ChainId:      uint(params.Order.DestChain.Uint64()),
		Swapica:      e.swapperContract.String(),
		OrderId:      uint(params.Order.Id.Uint64()),
		TokenToSell:  params.Order.TokenToSell.String(),
		AmountToSell: uint(params.Order.AmountToSell.Uint64()),
		OriginChain:  uint(srcChainParams.ChainId),
	})
	if err != nil {
		return nil, err
	}

	sender := common.HexToAddress(params.Sender)

	data := signature.OrderData{
		OrderData: calldata,
	}
	sign, err := e.signer.Sign(&data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign calldata")
	}

	tx, err := e.swapper.CreateMatch(buildTransactOpts(sender), calldata, [][]byte{sign})
	if err != nil {
		return nil, err
	}

	return encodeTx(tx, sender, params.Order.DestChain, params.DestChain.ID, nil)
}
