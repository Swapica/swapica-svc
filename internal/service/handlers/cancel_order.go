package handlers

import (
	"math/big"
	"net/http"

	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/Swapica/swapica-svc/internal/service/models"
	"github.com/Swapica/swapica-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCancelOrderRequest(r)
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

	order, err := ProxyRepo(r).Get(srcChain.ID).GetOrder(big.NewInt(int64(request.OrderId)))
	if err != nil {
		Log(r).WithError(err).Error("failed to get order")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	orderStatus, err := ProxyRepo(r).Get(srcChain.ID).GetOrderStatus(order.Id)
	if err != nil {
		Log(r).WithError(err).Error("failed to get order status")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tx, err := ProxyRepo(r).Get(request.SrcChain).CancelOrder(types.CancelOrderParams{
		Sender:      request.Sender,
		SrcChain:    *srcChain,
		Order:       order,
		OrderStatus: orderStatus,
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to create cancel order transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if tx == nil {
		Log(r).WithError(err).Error("failed to build transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, models.NewTxResponse(tx, *srcChain))
}
