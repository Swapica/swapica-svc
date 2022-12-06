package evm

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

func EncodeCancelMatch(selector uint, chainId *big.Int, swapica common.Address, id *big.Int) (string, error) {
	structThing, err := abi.NewType("tuple", "struct thing", []abi.ArgumentMarshaling{
		{Name: "selector", Type: "uint8"},
		{Name: "chain_id", Type: "uint"},
		{Name: "swapica", Type: "address"},
		{Name: "id", Type: "uint"},
	})
	if err != nil {
		return "", errors.New("failed to create abi structure")
	}

	args := abi.Arguments{
		{Type: structThing, Name: "param_one"},
	}

	record := struct {
		Selector uint
		ChainId  *big.Int
		Swapica  common.Address
		Id       *big.Int
	}{
		selector,
		chainId,
		swapica,
		id,
	}

	packed, err := args.Pack(&record)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(packed), err
}
