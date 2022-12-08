package evm

import (
	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
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

	return encodeTx(tx, sender, e.chainID, params.SrcChain.ID, nil)
}

func (e *evmProxy) executeMatchErc20(params types.ExecuteMatchParams, sender common.Address) (*ethTypes.Transaction, error) {
	orderData, err := EncodeExecuteMatch(executeMatchCalldata{
		Selector: executeMatch,
		ChainId:  uint(params.Order.DestChain.Uint64()),
		Swapica:  e.swapperContract.String(),
		MatchId:  uint(params.Match.Id.Uint64()),
		Receiver: params.Order.Account.String(),
	})

	if ok, err := e.validateExecuteMatchErc20(params, sender); !ok {
		return nil, err
	}

	data := signature.OrderData{
		OrderData: orderData,
	}
	sign, err := e.signer.Sign(&data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign order data")
	}

	tx, err := e.swapper.ExecuteMatch(
		buildTransactOpts(sender),
		orderData,
		[][]byte{sign},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tx")
	}

	return tx, nil
}

func (e *evmProxy) validateExecuteMatchErc20(params types.ExecuteMatchParams, sender common.Address) (bool, error) {
	if params.OrderStatus.State != awaitingMatch {
		return false, errors.New("cannot cancel a match if order is canceled or executed by matcher")
	}

	if params.MatchStatus.State != awaitingFinalization {
		return false, errors.New("cannot execute a match when it is not awaiting finalization")
	}

	if params.Receiver != params.Order.Account.String() {
		return false, errors.New("invalid receiver")
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
