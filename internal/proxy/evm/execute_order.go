package evm

import (
	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/state"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (e *evmProxy) ExecuteOrder(params types.ExecuteOrderParams) (interface{}, error) {
	sender := common.HexToAddress(params.Sender)

	tx, err := e.executeOrderErc20(params, sender)
	if err != nil {
		return nil, err
	}

	return encodeTx(tx, sender, e.chainID, params.SrcChain.ID, nil)
}

func (e *evmProxy) executeOrderErc20(params types.ExecuteOrderParams, sender common.Address) (*ethTypes.Transaction, error) {
	orderData, err := EncodeExecuteOrder(executeOrderCalldata{
		Selector: executeOrder,
		ChainId:  uint(params.Order.DestChain.Uint64()),
		Swapica:  e.swapperContract.String(),
		OrderId:  uint(params.Order.Id.Uint64()),
		Receiver: params.Match.Account.String(),
		MatchId:  uint(params.Match.Id.Uint64()),
	})

	if ok, err := e.validateExecuteOrderErc20(params, sender); !ok {
		return nil, err
	}

	data := signature.OrderData{
		OrderData: orderData,
	}
	sign, err := e.signer.Sign(&data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign order data")
	}

	tx, err := e.swapper.ExecuteOrder(
		buildTransactOpts(sender),
		orderData,
		[][]byte{sign},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tx")
	}

	return tx, nil
}

func (e *evmProxy) validateExecuteOrderErc20(params types.ExecuteOrderParams, sender common.Address) (bool, error) {
	if params.Match.Account != sender {
		return false, errors.New("invalid sender")
	}

	if params.OrderStatus.State != state.AwaitingMatch {
		return false, errors.New("cannot execute order if it is not awaiting match")
	}

	if params.MatchStatus.State != state.AwaitingFinalization {
		return false, errors.New("cannot execute order if match status is not awaiting finalization")
	}

	if params.Order.DestChain != params.Match.OriginChain {
		return false, errors.New("discrepancy between order and match chains")
	}

	if params.Order.AmountToBuy != params.Match.AmountToSell {
		return false, errors.New("mismatch between order amount to buy and match amount to sell")
	}

	if params.Order.TokenToBuy != params.Match.TokenToSell {
		return false, errors.New("mismatch between order token to buy and match token to sell")
	}

	if params.Order.Id != params.Match.OriginOrderId {
		return false, errors.New("mismatch between order id and match origin order id")
	}

	return true, nil
}
