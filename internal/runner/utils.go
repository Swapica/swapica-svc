package runner

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

type executeCalldata struct {
	Token      common.Address
	Commission *big.Int
	Receiver   common.Address
	CoreData   []byte
}

func EncodeExecuteParams(calldata executeCalldata) (string, error) {
	calldataType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "token", Type: "address"},
		{Name: "commission", Type: "uint256"},
		{Name: "receiver", Type: "address"},
		{Name: "core_data", Type: "bytes"},
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
