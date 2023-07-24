package handlers

import (
	"math"
	"math/big"
	"net/http"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/enums"
	"github.com/Swapica/swapica-svc/internal/service/models"
	"github.com/Swapica/swapica-svc/internal/service/requests"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetCommissionEstimate(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCommissionEstimateRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	gasPrice, err := ProxyRepo(r).Get(request.DestChain).SuggestGasPrice()
	if err != nil {
		Log(r).WithError(err).Error("failed to suggest gas price")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	lowerGas := big.NewInt(100000)
	upperGas := big.NewInt(150000)

	lowerGas.Mul(lowerGas, gasPrice)
	upperGas.Mul(upperGas, gasPrice)

	lowerGasFloat := new(big.Float).SetInt(lowerGas)
	lowerGasNative := ConvertAmount(lowerGasFloat, 18)

	upperGasFloat := new(big.Float).SetInt(upperGas)
	upperGasNative := ConvertAmount(upperGasFloat, 18)

	amountToBuyFloat, _ := new(big.Float).SetString(request.AmountToBuy)

	tokenChains, err := TokenChainsQ(r).FilterByChainID(request.DestChain).Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to select token chains")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	nativeSymbol := ""
	for _, tokenChain := range tokenChains {
		if tokenChain.ContractAddress == nil {
			t, err := TokensQ(r).FilterByID(tokenChain.TokenID).Get()
			if err != nil {
				Log(r).WithError(err).Error("failed to get token")
				ape.RenderErr(w, problems.InternalError())
				return
			}
			if t == nil {
				Log(r).Error("token not found")
				ape.RenderErr(w, problems.InternalError())
				return
			}
			nativeSymbol = t.Symbol
		}
	}

	if request.TokenToBuy == enums.TokenTypeNative {
		amountToBuyFloat = ConvertAmount(amountToBuyFloat, 18)

		lowerCommission := getPercent(lowerGasNative, amountToBuyFloat)
		upperCommission := getPercent(upperGasNative, amountToBuyFloat)

		lowerCommission.Mul(amountToBuyFloat, lowerCommission)
		upperCommission.Mul(amountToBuyFloat, upperCommission)

		lowerCommissionUsd, err := Converter(r).ConvertToStat(lowerCommission, nativeSymbol)
		if err != nil {
			Log(r).WithError(err).Error("failed to convert")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		upperCommissionUsd, err := Converter(r).ConvertToStat(upperCommission, nativeSymbol)
		if err != nil {
			Log(r).WithError(err).Error("failed to convert")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, models.NewCommissionEstimateResponse(lowerCommission, lowerCommissionUsd, upperCommission, upperCommissionUsd))
		return
	}

	instanceErc20, err := ProxyRepo(r).Get(request.DestChain).GetTokenInstance(common.HexToAddress(request.TokenToBuy))
	if err != nil {
		Log(r).WithError(err).Error("failed to get token instance")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	symbol, err := instanceErc20.Symbol(&bind.CallOpts{})
	if err != nil {
		Log(r).WithError(err).Error("failed to get token symbol")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	decimals, err := instanceErc20.Decimals(&bind.CallOpts{})
	if err != nil {
		Log(r).WithError(err).Error("failed to get token decimals")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	lowerCommInAsset, err := Converter(r).Convert(lowerGasNative, nativeSymbol, symbol)
	if err != nil {
		Log(r).WithError(err).Error("failed to convert")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	upperCommInAsset, err := Converter(r).Convert(upperGasNative, nativeSymbol, symbol)
	if err != nil {
		Log(r).WithError(err).Error("failed to convert")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	amountToBuyFloat = ConvertAmount(amountToBuyFloat, decimals)

	lowerCommission := getPercent(lowerCommInAsset, amountToBuyFloat)
	upperCommission := getPercent(upperCommInAsset, amountToBuyFloat)

	lowerCommission.Mul(amountToBuyFloat, lowerCommission)
	upperCommission.Mul(amountToBuyFloat, upperCommission)

	lowerCommissionUsd, err := Converter(r).ConvertToStat(lowerCommission, nativeSymbol)
	if err != nil {
		Log(r).WithError(err).Error("failed to convert")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	upperCommissionUsd, err := Converter(r).ConvertToStat(upperCommission, nativeSymbol)
	if err != nil {
		Log(r).WithError(err).Error("failed to convert")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, models.NewCommissionEstimateResponse(lowerCommission, lowerCommissionUsd, upperCommission, upperCommissionUsd))
}

func ConvertAmount(weiFloat *big.Float, decimals uint8) *big.Float {
	ether := new(big.Float)
	decimal := new(big.Float).SetInt(big.NewInt(int64(math.Pow10(int(decimals)))))
	ether.Quo(weiFloat, decimal)

	return ether
}

func getPercent(x *big.Float, y *big.Float) *big.Float {
	multiply := x.Mul(x, big.NewFloat(100))
	return multiply.Quo(multiply, y)
}
