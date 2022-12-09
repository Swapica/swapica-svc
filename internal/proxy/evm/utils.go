package evm

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/Swapica/swapica-svc/resources"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

const gasLimit = 300000

func skipSig(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
	return transaction, nil
}

func buildTransactOpts(from common.Address) *bind.TransactOpts {
	return buildTransactOptsWithValue(from, nil)
}

func buildTransactOptsWithValue(from common.Address, value *big.Int) *bind.TransactOpts {
	return &bind.TransactOpts{
		From:     from,
		Signer:   skipSig,
		NoSend:   true,
		GasLimit: gasLimit,
		Value:    value,
	}
}

func parseChainParams(params json.RawMessage) (resources.EvmChainParams, error) {
	chainParams := resources.EvmChainParams{}
	err := json.Unmarshal(params, &chainParams)
	return chainParams, err
}

func encodeTx(tx *types.Transaction, from common.Address, chainID *big.Int, chain string, confirmed *bool) (interface{}, error) {
	return resources.EvmTransaction{
		Key: resources.Key{
			ID:   tx.Hash().String(),
			Type: resources.EVM_TRANSACTION,
		},
		Attributes: resources.EvmTransactionAttributes{
			Confirmed: confirmed,
			TxBody: resources.EvmTransactionTxBody{
				From:    from.String(),
				To:      tx.To().String(),
				Value:   tx.Value().String(),
				Data:    hexutil.Encode(tx.Data()),
				ChainId: fmt.Sprintf("0x%x", chainID.Int64()),
			},
		},
		Relationships: resources.EvmTransactionRelationships{
			Chain: resources.Key{
				ID:   chain,
				Type: resources.CHAIN,
			}.AsRelation(),
		},
	}, nil
}
