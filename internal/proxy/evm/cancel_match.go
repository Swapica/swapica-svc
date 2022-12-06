package evm

import (
	"github.com/Swapica/swapica-svc/internal/proxy"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"math/big"
)

func (e *evmProxy) CancelMatch(params types.CancelMatchParams, repo proxy.ProxyRepo) (interface{}, error) {
	tokenType := TokenTypeErc20

	sender := common.HexToAddress(params.Sender)

	var tx *ethTypes.Transaction
	var err error

	switch tokenType {
	case TokenTypeErc20:
		tx, err = e.cancelMatchErc20(params, sender, repo)
	default:
		return nil, errors.New("unsupported token type")
	}

	if err != nil {
		return nil, err
	}
	if tx == nil {
		// Token is already approved
		return nil, nil
	}

	return encodeTx(tx, sender, e.chainID, params.SrcChain.ID, nil)
}

func (e *evmProxy) cancelMatchErc20(params types.CancelMatchParams, sender common.Address, repo proxy.ProxyRepo) (*ethTypes.Transaction, error) {
	matchId := new(big.Int).SetInt64(int64(params.MatchId))

	match, err := e.swapper.Matches(&bind.CallOpts{}, matchId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check match")
	}
	if match.Account != sender {
		return nil, errors.Wrap(err, "invalid sender")
	}

	order, err := repo.Get(params.OrgChain.ID).GetOrder(match.OriginOrderId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check order")
	}
	if order.DestChain != match.OriginChain || order.AmountToBuy != match.AmountToSell ||
		order.TokenToBuy != match.TokenToSell || matchId != match.Id {
		return nil, errors.Wrap(err, "ivalid options provided by user")
	}

	orderData, err := EncodeCancelMatch(
		3,
		order.DestChain,
		e.swapperContract,
		matchId,
	)

	log := signature.CancelMatch{
		OrderData: orderData,
	}
	sign, err := e.signer.Sign(&log)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign log")
	}

	tx, err := e.swapper.CancelMatch(
		buildTransactOpts(sender),
		[]byte(orderData),
		append([][]byte{}, sign),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tx")
	}

	return tx, nil
}
