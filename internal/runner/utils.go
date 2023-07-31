package runner

import (
	"context"
	"math"
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/enums"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/erc20"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var CommissionIsTooHigh error = errors.New("commission is >= 100%")

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

func (r *Runner) CommissionEstimate(
	txData []byte, contractAddress common.Address, tokenAddress common.Address, nativeSymbol, amount, rpc string,
) (*big.Int, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: txData,
	}

	gasLimit, err := client.EstimateGas(context.Background(), callMsg)
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	gasLimitBig := big.NewInt(int64(gasLimit))
	gas := gasLimitBig.Mul(gasLimitBig, gasPrice)

	gasInNative := ConvertAmount(gas, 18)

	amountToInt, _ := new(big.Int).SetString(amount, 10)

	if tokenAddress == common.HexToAddress(enums.TokenTypeNative) {
		commission := getPercent(gasInNative, ConvertAmount(amountToInt, 18))

		return ConvertToBigIntCommission(commission), nil
	}

	instanceErc20, err := erc20.NewErc20(contractAddress, client)
	if err != nil {
		return nil, err
	}

	symbol, err := instanceErc20.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	decimals, err := instanceErc20.Decimals(&bind.CallOpts{})

	commissionInAsset, err := r.converter.Convert(gasInNative, nativeSymbol, symbol)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert")
	}

	commissionPercent := getPercent(commissionInAsset, ConvertAmount(amountToInt, decimals))

	if commissionPercent.Cmp(big.NewInt(100)) > -1 {
		return commissionPercent, CommissionIsTooHigh
	}

	return ConvertToBigIntCommission(commissionPercent), nil
}

func ConvertToBigIntCommission(x *big.Int) *big.Int {
	base := big.NewInt(10)
	exponent := big.NewInt(27)
	result := new(big.Int).Exp(base, exponent, nil)
	temp := new(big.Int).Mul(x, result)

	return new(big.Int).Div(temp, big.NewInt(100))
}

func ConvertAmount(wei *big.Int, decimals uint8) *big.Float {
	ether := new(big.Float)
	weiFloat := new(big.Float).SetInt(wei)
	decimal := new(big.Float).SetInt(big.NewInt(int64(math.Pow10(int(decimals)))))
	ether.Quo(weiFloat, decimal)

	return ether
}

func getPercent(x *big.Float, y *big.Float) *big.Int {
	multiply := x.Mul(x, big.NewFloat(100))
	divide := multiply.Quo(multiply, y)
	float64Value, _ := divide.Float64()
	ceilValue := math.Ceil(float64Value) // TODO remove

	return big.NewInt(int64(ceilValue))
}
