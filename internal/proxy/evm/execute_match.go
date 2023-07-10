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

func (e *evmProxy) ExecuteMatch(params types.ExecuteMatchParams) (interface{}, error) {
	sender := common.HexToAddress(params.Sender)

	threshold, err := e.getThreshold()
	if err != nil {
		return nil, err
	}

	tx, err := e.executeMatch(params, sender)
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

func (e *evmProxy) executeMatch(params types.ExecuteMatchParams, sender common.Address) (*ethTypes.Transaction, error) {
	orderData, err := EncodeExecuteMatch(executeMatchCalldata{
		Selector: executeMatch,
		ChainId:  params.Order.DestinationChain,
		Swapica:  e.swapperContract,
		MatchId:  params.Match.MatchId,
		Receiver: common.HexToAddress(params.Receiver),
	})
	if err != nil {
		return nil, err
	}

	if ok, err := e.validateExecuteMatch(params); !ok {
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

func (e *evmProxy) validateExecuteMatch(params types.ExecuteMatchParams) (bool, error) {
	if enums.State(params.Order.Status.State) != enums.Executed {
		return false, errors.New("cannot execute a match if order is not executed")
	}

	if params.Order.Status.MatchId.String() != params.Match.MatchId.String() {
		return false, errors.New("cannot execute a match if order executed by someone other match")
	}

	if params.Order.Status.MatchSwapica.String() != e.swapperContract.String() {
		return false, errors.New("cannot execute a match if the match is not created by this swapica contract")
	}

	if enums.State(params.Match.State) != enums.AwaitingFinalization {
		return false, errors.New("cannot execute a match when it is not awaiting finalization")
	}

	if params.Receiver != params.Order.Creator.String() && params.Receiver != e.relayerContract.String() {
		return false, errors.New("invalid receiver")
	}

	if params.Order.AmountToBuy.String() != params.Match.AmountToSell.String() {
		return false, errors.New("mismatch between order amount to buy and match amount to sell")
	}

	if params.Order.TokenToBuy != params.Match.TokenToSell {
		return false, errors.New("mismatch between order token to buy and match token to sell")
	}

	if params.Order.OrderId.String() != params.Match.OriginOrderId.String() {
		return false, errors.New("mismatch between order id and match origin order id")
	}

	return true, nil
}
