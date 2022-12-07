package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Selector uint8

const (
	executeOrder Selector = iota
	executeMatch
	createMatch
	cancelMatch
)

type createMatchCalldata struct {
	Selector
	ChainId      uint
	Swapica      string
	OrderId      uint
	TokenToSell  string
	AmountToSell uint
	OriginChain  uint
}

func CreateMatchCalldata(calldata createMatchCalldata) ([]byte, error) {
	calldataType, err := abi.NewType("turple", "", []abi.ArgumentMarshaling{
		{Name: "selector", Type: "uint8"},
		{Name: "chain_id", Type: "uint"},
		{Name: "swapica", Type: "address"},
		{Name: "order_id", Type: "uint"},
		{Name: "token_to_sale", Type: "address"},
		{Name: "amount_to_sell", Type: "uint"},
		{Name: "origin_chain", Type: "uint"},
	})
	if err != nil {
		return nil, err
	}

	calldataArgs := abi.Arguments{
		{Type: calldataType, Name: "calldata"},
	}

	packed, err := calldataArgs.Pack(&calldata)
	if err != nil {
		return nil, err
	}

	return []byte(hexutil.Encode(packed)), nil
}

type cancelMatchCalldata struct {
	Selector
	ChainId uint
	Swapica string
	MatchId uint
}

func EncodeCancelMatch(calldata cancelMatchCalldata) ([]byte, error) {
	calldataType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "selector", Type: "uint8"},
		{Name: "chain_id", Type: "uint"},
		{Name: "swapica", Type: "address"},
		{Name: "id", Type: "uint"},
	})
	if err != nil {
		return nil, err
	}

	args := abi.Arguments{
		{Type: calldataType, Name: "calldata"},
	}

	packed, err := args.Pack(&calldata)
	if err != nil {
		return nil, err
	}
	return []byte(hexutil.Encode(packed)), err
}
