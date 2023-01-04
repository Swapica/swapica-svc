package evm

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/enums"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/erc1155"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/erc20"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/erc721"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (e *evmProxy) Approve(params types.ApproveParams) (interface{}, error) {
	fromAddress := common.HexToAddress(params.Sender)

	var tx *ethTypes.Transaction
	var err error
	switch params.TokenType {
	case enums.TokenTypeNative:
		// ApproveParams not needed for native token
		return nil, nil
	case enums.TokenTypeErc20:
		tx, err = e.approveErc20(common.HexToAddress(params.TokenAddress), fromAddress)
	case enums.TokenTypeErc721:
		tx, err = e.approveErc721(common.HexToAddress(params.TokenAddress), fromAddress)
	case enums.TokenTypeErc1155:
		tx, err = e.approveErc1155(common.HexToAddress(params.TokenAddress), fromAddress)
	default:
		return nil, errors.New("unknown token type")
	}
	if err != nil {
		return nil, err
	}
	if tx == nil {
		// Token is already approved
		return nil, nil
	}

	confirmed := true
	return encodeTx(tx, fromAddress, e.chainID, params.Chain.ID, &confirmed)
}

func (e *evmProxy) approveErc20(tokenAddress common.Address, approveFrom common.Address) (*ethTypes.Transaction, error) {
	token, err := erc20.NewErc20(tokenAddress, e.client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transactor")
	}

	amount, err := token.Allowance(&bind.CallOpts{}, approveFrom, e.swapperContract)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check allowance")
	}
	if amount.Cmp(big.NewInt(0)) == 1 {
		// Token is already approved
		return nil, nil
	}

	uin256max := big.NewInt(0)
	uin256max.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	tx, err := token.Approve(buildTransactOpts(approveFrom), e.swapperContract, uin256max)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build token approve tx")
	}

	return tx, nil
}

func (e *evmProxy) approveErc721(tokenAddress common.Address, approveFrom common.Address) (*ethTypes.Transaction, error) {
	token, err := erc721.NewErc721(tokenAddress, e.client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transactor")
	}

	approved, err := token.IsApprovedForAll(&bind.CallOpts{}, approveFrom, e.swapperContract)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check approval")
	}
	if approved {
		// Token is already approved
		return nil, nil
	}

	tx, err := token.SetApprovalForAll(buildTransactOpts(approveFrom), e.swapperContract, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build token approve tx")
	}

	return tx, err
}

func (e *evmProxy) approveErc1155(tokenAddress common.Address, approveFrom common.Address) (*ethTypes.Transaction, error) {
	token, err := erc1155.NewErc1155(tokenAddress, e.client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transactor")
	}

	approved, err := token.IsApprovedForAll(&bind.CallOpts{}, approveFrom, e.swapperContract)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check approval")
	}
	if approved {
		// Token is already approved
		return nil, nil
	}

	tx, err := token.SetApprovalForAll(buildTransactOpts(approveFrom), e.swapperContract, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build token approve tx")
	}

	return tx, err
}
