package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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
	ChainId      *big.Int
	Swapica      common.Address
	OrderId      *big.Int
	TokenToSell  common.Address
	AmountToSell *big.Int
	OriginChain  *big.Int
}

func CreateMatchCalldata(calldata createMatchCalldata) (string, error) {
	calldataType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "selector", Type: "uint8"},
		{Name: "chain_id", Type: "uint256"},
		{Name: "swapica", Type: "address"},
		{Name: "order_id", Type: "uint256"},
		{Name: "token_to_sell", Type: "address"},
		{Name: "amount_to_sell", Type: "uint256"},
		{Name: "origin_chain", Type: "uint256"},
	})
	if err != nil {
		return "", err
	}

	calldataArgs := abi.Arguments{
		{Type: calldataType, Name: "calldata"},
	}

	packed, err := calldataArgs.Pack(&calldata)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(packed), nil
}

type cancelMatchCalldata struct {
	Selector
	ChainId *big.Int
	Swapica common.Address
	MatchId *big.Int
}

func EncodeCancelMatch(calldata cancelMatchCalldata) (string, error) {
	calldataType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "selector", Type: "uint8"},
		{Name: "chain_id", Type: "uint256"},
		{Name: "swapica", Type: "address"},
		{Name: "match_id", Type: "uint256"},
	})
	if err != nil {
		return "", err
	}

	args := abi.Arguments{
		{Type: calldataType, Name: "calldata"},
	}

	packed, err := args.Pack(&calldata)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(packed), nil
}

type executeMatchCalldata struct {
	Selector
	ChainId  *big.Int
	Swapica  common.Address
	MatchId  *big.Int
	Receiver common.Address
}

func EncodeExecuteMatch(calldata executeMatchCalldata) (string, error) {
	calldataType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "selector", Type: "uint8"},
		{Name: "chain_id", Type: "uint256"},
		{Name: "swapica", Type: "address"},
		{Name: "match_id", Type: "uint256"},
		{Name: "receiver", Type: "address"},
	})
	if err != nil {
		return "", err
	}

	args := abi.Arguments{
		{Type: calldataType, Name: "calldata"},
	}

	packed, err := args.Pack(&calldata)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(packed), err
}

type executeOrderCalldata struct {
	Selector
	ChainId      *big.Int
	Swapica      common.Address
	OrderId      *big.Int
	Receiver     common.Address
	MatchSwapica common.Address
	MatchId      *big.Int
}

func EncodeExecuteOrder(calldata executeOrderCalldata) (string, error) {
	calldataType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "selector", Type: "uint8"},
		{Name: "chain_id", Type: "uint256"},
		{Name: "swapica", Type: "address"},
		{Name: "order_id", Type: "uint256"},
		{Name: "receiver", Type: "address"},
		{Name: "match_swapica", Type: "address"},
		{Name: "match_id", Type: "uint256"},
	})
	if err != nil {
		return "", err
	}

	args := abi.Arguments{
		{Type: calldataType, Name: "calldata"},
	}

	packed, err := args.Pack(&calldata)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(packed), err
}
