package evm

import (
	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (e *evmProxy) CancelMatch(params types.CancelMatchParams) (interface{}, error) {
	sender := common.HexToAddress(params.Sender)

	tx, err := e.cancelMatchErc20(params, sender)
	if err != nil {
		return nil, err
	}

	return encodeTx(tx, sender, e.chainID, params.SrcChain.ID, nil)
}

func (e *evmProxy) cancelMatchErc20(params types.CancelMatchParams, sender common.Address) (*ethTypes.Transaction, error) {
	orderData, err := EncodeCancelMatch(cancelMatchCalldata{
		Selector: cancelMatch,
		ChainId:  uint(params.Order.DestChain.Uint64()),
		Swapica:  e.swapperContract.String(),
		MatchId:  uint(params.Match.Id.Uint64()),
	})

	if ok, err := e.validateCancelMatchErc20(params, sender); !ok {
		return nil, err
	}

	data := signature.OrderData{
		OrderData: orderData,
	}
	sign, err := e.signer.Sign(&data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign order data")
	}

	tx, err := e.swapper.CancelMatch(
		buildTransactOpts(sender),
		orderData,
		[][]byte{sign},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tx")
	}

	return tx, nil
}

func (e *evmProxy) validateCancelMatchErc20(params types.CancelMatchParams, sender common.Address) (bool, error) {
	if params.Match.Account != sender {
		return false, errors.New("invalid sender")
	}

	if params.OrderStatus.State != canceled ||
		(params.OrderStatus.State == executed && params.OrderStatus.ExecutedBy == params.Match.Id) {
		return false, errors.New("cannot cancel a match if order is canceled or executed by matcher")
	}

	if params.MatchStatus.State != awaitingFinalization {
		return false, errors.New("cannot cancel a match when it is not awaiting finalization")
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
