package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/erc20"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"io"
	"math"
	"math/big"
	"net/http"
	"strings"
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

const (
	nativeToken = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
)

func CommissionEstimate(txData []byte, contractAddress common.Address, tokenAddress common.Address, amount string, rpc string) (*big.Int, error) {
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

	if tokenAddress == common.HexToAddress(nativeToken) {
		return getPercent(gasInNative, ConvertAmount(amountToInt, 18)), nil
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
	price, err := GetTokenPrice(symbol)

	return getPercent(gasInNative.Quo(gasInNative, new(big.Float).SetFloat64(price)), ConvertAmount(amountToInt, decimals)), nil
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
	ceilValue := math.Ceil(float64Value)

	return big.NewInt(int64(ceilValue))
}

func GetTokenPrice(symbol string) (float64, error) {
	if strings.ToLower(symbol) == "usdc" || strings.ToLower(symbol) == "usdt" {
		symbol = "usd"
	}

	resp, err := http.Get(fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=eth", symbol))
	if err != nil {
		return 0, errors.Wrap(err, "failed to get price")
	}

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.Wrap(err, "failed to read response")
	}

	if err != nil {
		return 0, errors.Wrap(err, "failed to unmarshal response")
	}

	var data map[string]map[string]float64
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return 0, errors.Wrap(err, "failed to unmarshal response")
	}

	price := data[strings.ToLower(symbol)]["eth"]

	return price, nil
}
