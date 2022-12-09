package evm

import (
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (e *evmProxy) CancelOrder(params types.CancelOrderParams) (interface{}, error) {
	sender := common.HexToAddress(params.Sender)

	tx, err := e.cancelOrderErc20(params, sender)
	if err != nil {
		return nil, err
	}

	return encodeTx(tx, sender, e.chainID, params.SrcChain.ID, nil)
}

func (e *evmProxy) cancelOrderErc20(params types.CancelOrderParams, sender common.Address) (*ethTypes.Transaction, error) {
	if params.OrderStatus.State != awaitingMatch {
		return nil, errors.New("cam mot cancel order if it is not awaiting match")
	}

	tx, err := e.swapper.CancelOrder(
		buildTransactOpts(sender),
		params.Order.Id,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cancel tx")
	}

	return tx, nil
}
