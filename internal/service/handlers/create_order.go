package handlers

import (
	"errors"
	"math/big"
	"net/http"

	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/Swapica/swapica-svc/internal/service/models"
	"github.com/Swapica/swapica-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateOrderRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	srcChain, err := ChainsQ(r).FilterByID(request.SrcChain).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get source chain")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if srcChain == nil {
		Log(r).Error("source chain not found")
		ape.RenderErr(w, problems.BadRequest(errors.New("source chain not found"))...)
		return
	}

	destChain, err := ChainsQ(r).FilterByID(request.DestChain).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get destination chain")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if destChain == nil {
		Log(r).Error("destination chain not found")
		ape.RenderErr(w, problems.BadRequest(errors.New("destination chain not found"))...)
		return
	}

	amountToBuy,ok1:=big.NewInt(0).SetString(request.AmountToBuy,10)
	amountToSell,ok2:=big.NewInt(0).SetString(request.AmountToSell,10)
	if !ok1||!ok2{
		msg:="amount_to_buy or amount_to_sell is not big integer"
		Log(r).WithFields(logan.F{"amount_to_buy": request.AmountToBuy,"amount_to_sell": request.AmountToSell}).Debug(msg)
		ape.RenderErr(w, problems.BadRequest(nil)...)
		return
	}

	tx, err := ProxyRepo(r).Get(request.SrcChain).CreateOrder(types.CreateOrderParams{
		Sender:       request.Sender,
		SrcChain:     *srcChain,
		TokenToSell:  request.TokenToSell,
		AmountToSell: amountToSell,
		TokenToBuy:   request.TokenToBuy,
		AmountToBuy:  amountToBuy,
		DestChain:    *destChain,
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to make create order transaction")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	if tx == nil {
		Log(r).WithError(err).Error("failed to build transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, models.NewTxResponse(tx, *srcChain))
}
