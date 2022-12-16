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

func (e *evmProxy) ExecuteOrder(params types.ExecuteOrderParams) (interface{}, error) {
	sender := common.HexToAddress(params.Sender)

	tx, err := e.executeOrderErc20(params, sender)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	return encodeTx(tx, sender, e.chainID, params.SrcChain.ID, nil)
}

func (e *evmProxy) executeOrderErc20(params types.ExecuteOrderParams, sender common.Address) (*ethTypes.Transaction, error) {
	orderData, err := EncodeExecuteOrder(executeOrderCalldata{
		Selector: executeOrder,
		ChainId:  params.Match.OriginChain,
		Swapica:  e.swapperContract,
		OrderId:  params.Order.Id,
		Receiver: params.Match.Account,
		MatchId:  params.Match.Id,
	})
	if err != nil {
		return nil, err
	}

	if ok, err := e.validateExecuteOrderErc20(params); !ok {
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

	tx, err := e.swapper.ExecuteOrder(
		buildTransactOpts(sender),
		hexedCalldata,
		append([][]byte{}, sign),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tx")
	}

	return tx, nil
}

func (e *evmProxy) validateExecuteOrderErc20(params types.ExecuteOrderParams) (bool, error) {
	if params.Receiver != params.Match.Account.String() {
		return false, errors.New("invalid receiver")
	}

	if params.OrderStatus.State != enums.AwaitingMatch {
		return false, errors.New("cannot execute order if it is not awaiting match")
	}

	if params.MatchStatus.State != enums.AwaitingFinalization {
		return false, errors.New("cannot execute order if match status is not awaiting finalization")
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
