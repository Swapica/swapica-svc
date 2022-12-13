package evm

import (
	"encoding/hex"
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/state"
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
		ChainId:      params.Order.DestChain,
		Swapica:      e.swapperContract,
		OrderId:      params.Order.Id,
		TokenToSell:  params.Order.TokenToBuy,
		AmountToSell: params.Order.AmountToBuy,
		OriginChain:  big.NewInt(srcChainParams.ChainId),
	})
	if err != nil {
		return nil, err
	}

	if params.OrderStatus.State != state.AwaitingMatch {
		return nil, errors.New("can not create match if order is not awaiting match")
	}

	sender := common.HexToAddress(params.Sender)

	data := signature.OrderData{
		OrderData: calldata,
	}
	sign, err := e.signer.Sign(&data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign calldata")
	}
	hexedCalldata, err := hex.DecodeString(calldata[2:])
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode calldata")
	}

	tx, err := e.swapper.CreateMatch(
		GetTransactionOpts(params.Order.TokenToBuy.String(),
			sender, params.Order.AmountToBuy),
		hexedCalldata,
		append([][]byte{}, sign),
	)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	return encodeTx(tx, sender, params.Order.DestChain, params.DestChain.ID, nil)
}
