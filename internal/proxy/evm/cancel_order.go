package evm

import (
	"github.com/Swapica/swapica-svc/internal/proxy/evm/enums"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (e *evmProxy) CancelOrder(params types.CancelOrderParams) (interface{}, error) {
	sender := common.HexToAddress(params.Sender)

	tx, err := e.cancelOrder(params, sender)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	confirmed := true
	return encodeTx(tx, sender, e.chainID, params.SrcChain.ID, &confirmed)
}

func (e *evmProxy) cancelOrder(params types.CancelOrderParams, sender common.Address) (*ethTypes.Transaction, error) {
	if enums.State(params.Order.Status.State) != enums.AwaitingMatch {
		return nil, errors.New("can not cancel order if it is not awaiting match")
	}

	tx, err := e.swapper.CancelOrder(
		buildTransactOpts(sender),
		params.Order.OrderId,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cancel tx")
	}

	return tx, nil
}
