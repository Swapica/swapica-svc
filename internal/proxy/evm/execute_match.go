package evm

import (
	"encoding/hex"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/state"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (e *evmProxy) ExecuteMatch(params types.ExecuteMatchParams) (interface{}, error) {
	sender := common.HexToAddress(params.Sender)

	tx, err := e.executeMatchErc20(params, sender)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	return encodeTx(tx, sender, e.chainID, params.SrcChain.ID, nil)
}

func (e *evmProxy) executeMatchErc20(params types.ExecuteMatchParams, sender common.Address) (*ethTypes.Transaction, error) {
	orderData, err := EncodeExecuteMatch(executeMatchCalldata{
		Selector: executeMatch,
		ChainId:  params.Order.DestChain,
		Swapica:  e.swapperContract,
		MatchId:  params.Match.Id,
		Receiver: params.Order.Account,
	})

	if ok, err := e.validateExecuteMatchErc20(params); !ok {
		return nil, err
	}

	data := signature.OrderData{
		OrderData: orderData,
	}
	sign, err := e.signer.Sign(&data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign order data")
	}
	hexedCalldata, err := hex.DecodeString(orderData[2:])
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode calldata")
	}

	tx, err := e.swapper.ExecuteMatch(
		buildTransactOpts(sender),
		hexedCalldata,
		append([][]byte{}, sign),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tx")
	}

	return tx, nil
}

func (e *evmProxy) validateExecuteMatchErc20(params types.ExecuteMatchParams) (bool, error) {
	if params.OrderStatus.State != state.Executed {
		return false, errors.New("cannot execute a match if order is not executed")
	}

	if params.OrderStatus.ExecutedBy.String() != params.Match.Id.String() {
		return false, errors.New("cannot execute a match if order executed by someone other match")
	}

	if params.MatchStatus.State != state.AwaitingFinalization {
		return false, errors.New("cannot execute a match when it is not awaiting finalization")
	}

	if params.Receiver != params.Order.Account.String() {
		return false, errors.New("invalid receiver")
	}

	if params.Order.AmountToBuy.String() != params.Match.AmountToSell.String() {
		return false, errors.New("mismatch between order amount to buy and match amount to sell")
	}

	if params.Order.TokenToBuy != params.Match.TokenToSell {
		return false, errors.New("mismatch between order token to buy and match token to sell")
	}

	if params.Order.Id.String() != params.Match.OriginOrderId.String() {
		return false, errors.New("mismatch between order id and match origin order id")
	}

	return true, nil
}
