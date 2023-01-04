package handlers

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"net/http"

	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/Swapica/swapica-svc/internal/service/models"
	"github.com/Swapica/swapica-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func ExecuteOrder(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewExecuteOrderRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	destChain, err := ChainsQ(r).FilterByID(request.DestChain).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get destination chain")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if destChain == nil {
		Log(r).Debug("destination chain not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	srcChain, err := ChainsQ(r).FilterByID(request.SrcChain).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get source chain")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if srcChain == nil {
		Log(r).Debug("source chain not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	order, err := ProxyRepo(r).Get(request.SrcChain).GetOrder(big.NewInt(int64(request.OrderId)))
	if err != nil {
		Log(r).WithError(err).Error("failed to get order")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	match, err := ProxyRepo(r).Get(request.DestChain).GetMatch(big.NewInt(int64(request.MatchId)))
	if err != nil {
		Log(r).WithError(err).Error("failed to get match")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	matchStatus, err := ProxyRepo(r).Get(destChain.ID).GetMatchStatus(match.Id)
	if err != nil {
		Log(r).WithError(err).Error("failed to get match status")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	orderStatus, err := ProxyRepo(r).Get(srcChain.ID).GetOrderStatus(order.Id)
	if err != nil {
		Log(r).WithError(err).Error("failed to get order status")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	params := types.ExecuteOrderParams{
		SrcChain:    *srcChain,
		DestChain:   *destChain,
		Order:       order,
		Match:       match,
		OrderStatus: orderStatus,
		MatchStatus: matchStatus,
		Receiver:    request.Receiver,
		Sender:      request.Sender,
	}

	if request.RawTxData != nil {
		rawTxData, err := hexutil.Decode(*request.RawTxData)
		if err != nil {
			return
		}

		params.RawTxData = &rawTxData
	}

	tx, err := ProxyRepo(r).Get(request.SrcChain).ExecuteOrder(params)
	if err != nil {
		Log(r).WithError(err).Error("failed to create match")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	if tx == nil {
		Log(r).WithError(err).Error("failed to build transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, models.NewTxResponse(tx, *destChain))
}
