package handlers

import (
	"errors"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/Swapica/swapica-svc/internal/service/models"
	"github.com/Swapica/swapica-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
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

	tx, err := ProxyRepo(r).Get(request.SrcChain).CreateOrder(types.CreateOrderParams{
		Sender:       request.Sender,
		SrcChain:     *srcChain,
		TokenToSell:  request.TokenToSell,
		AmountToSell: request.AmountToSell,
		TokenToBuy:   request.TokenToBuy,
		AmountToBuy:  request.AmountToBuy,
		DestChain:    *destChain,
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to make create order transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, models.NewTxResponse(tx, *srcChain))
}
