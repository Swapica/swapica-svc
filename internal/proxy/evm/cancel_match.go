package evm

import (
	"github.com/Swapica/swapica-svc/internal/proxy/evm/enums"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (e *evmProxy) CancelMatch(params types.CancelMatchParams) (interface{}, error) {
	sender := common.HexToAddress(params.Sender)

	threshold, err := e.getThreshold()
	if err != nil {
		return nil, err
	}

	tx, err := e.cancelMatch(params, sender)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	signNumber := int64(1)

	// if tx provided check it and sign; otherwise use created tx
	if params.RawTxData != nil {
		tx, signNumber, err = e.checkTxDataAndSign(
			buildTransactOpts(sender),
			tx,
			*params.RawTxData,
		)
		if err != nil {
			return nil, err
		}
	}

	confirmed := signNumber >= threshold

	return encodeTx(tx, sender, e.chainID, params.SrcChain.ID, &confirmed)
}

func (e *evmProxy) cancelMatch(params types.CancelMatchParams, sender common.Address) (*ethTypes.Transaction, error) {
	orderData, err := EncodeCancelMatch(cancelMatchCalldata{
		Selector: cancelMatch,
		ChainId:  params.Order.DestChain,
		Swapica:  e.swapperContract,
		MatchId:  params.Match.Id,
	})
	if err != nil {
		return nil, err
	}

	if ok, err := e.validateCancelMatch(params, sender); !ok {
		return nil, err
	}

	data := signature.OrderData{
		OrderData: orderData,
	}
	sign, err := e.signer.Sign(&data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign order data")
	}
	hexedCalldata, err := hexutil.Decode(orderData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode calldata")
	}

	tx, err := e.swapper.CancelMatch(
		buildTransactOpts(sender),
		hexedCalldata,
		append([][]byte{}, sign),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tx")
	}

	return tx, nil
}

func (e *evmProxy) validateCancelMatch(params types.CancelMatchParams, sender common.Address) (bool, error) {
	if params.Match.Account != sender {
		return false, errors.New("invalid sender")
	}

	notCanceled := enums.State(params.OrderStatus.State) != enums.Canceled
	isExecutedBy := enums.State(params.OrderStatus.State) == enums.Executed && params.OrderStatus.ExecutedBy == params.Match.Id
	isExecutedByThis := params.OrderStatus.MatchSwapica == e.swapperContract
	if notCanceled && isExecutedBy && isExecutedByThis {
		return false, errors.New("cannot cancel a match if order is canceled or executed by matcher")
	}

	if enums.State(params.MatchStatus.State) != enums.AwaitingFinalization {
		return false, errors.New("cannot cancel a match when it is not awaiting finalization")
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
