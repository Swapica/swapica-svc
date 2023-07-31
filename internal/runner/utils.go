package runner

import (
	"context"
	"math"
	"math/big"

	"github.com/ALTree/bigfloat"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/enums"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/erc20"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
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
	contractAddress common.Address, tokenAddress common.Address, nativeSymbol, amount, rpc string,
) (*big.Int, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	gasLimit := big.NewInt(150000)
	gas := gasLimit.Mul(gasLimit, gasPrice)

	gasInNative := ConvertAmount(gas, 18)

	amountToInt, _ := new(big.Int).SetString(amount, 10)

	if tokenAddress == common.HexToAddress(enums.TokenTypeNative) {
		commission, err := getPercent(gasInNative, ConvertAmount(amountToInt, 18))

		if err != nil {
			r.log.WithFields(logan.F{
				"commission_percent": commission.String(),
				"gas_limit":          150000,
				"gas_price":          gasPrice,
				"gas_in_native":      gasInNative.String(),
				"amount":             ConvertAmount(amountToInt, 18).String(),
			}).Warn("commission is too high")
			return commission, CommissionIsTooHigh
		}

		return commission, nil
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

	commissionPercent, err := getPercent(commissionInAsset, ConvertAmount(amountToInt, decimals))
	if err != nil {
		r.log.WithFields(logan.F{
			"commission_percent": commissionPercent.String(),
			"gas_limit":          150000,
			"gas_price":          gasPrice,
			"gas_in_native":      gasInNative.String(),
			"amount":             ConvertAmount(amountToInt, decimals).String(),
		}).Warn("commission is too high")
		return commissionPercent, CommissionIsTooHigh
	}

	return commissionPercent, nil
}

func ConvertAmount(wei *big.Int, decimals uint8) *big.Float {
	ether := new(big.Float)
	weiFloat := new(big.Float).SetInt(wei)
	decimal := new(big.Float).SetInt(big.NewInt(int64(math.Pow10(int(decimals)))))
	ether.Quo(weiFloat, decimal)

	return ether
}

func getPercent(x *big.Float, y *big.Float) (*big.Int, error) {
	res := new(big.Float)
	res.Quo(x, y)

	if res.Cmp(big.NewFloat(1)) > 0 {
		return nil, CommissionIsTooHigh
	}

	multiplier := bigfloat.Pow(big.NewFloat(10), big.NewFloat(27))
	res.Mul(res, multiplier)

	resI := new(big.Int)
	res.Int(resI)

	return resI, nil
}
