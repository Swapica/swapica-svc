package handlers

import (
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/Swapica/swapica-svc/internal/service/models"
	"github.com/Swapica/swapica-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"math/big"
	"net/http"
)

func ExecuteMatch(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewExecuteMatchRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	destChain, err := ChainsQ(r).FilterByID(request.DestChain).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get chains")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if destChain == nil {
		Log(r).Debug("chain not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	srcChain, err := ChainsQ(r).FilterByID(request.SrcChain).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get chains")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if srcChain == nil {
		Log(r).Debug("chain not found")
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

	matchStatus, err := ProxyRepo(r).Get(destChain.ID).GetMatchStatus(big.NewInt(int64(request.MatchId)))
	if err != nil {
		Log(r).WithError(err).Error("failed to get order")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	orderStatus, err := ProxyRepo(r).Get(srcChain.ID).GetOrderStatus(match.OriginOrderId)
	if err != nil {
		Log(r).WithError(err).Error("failed to get order")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tx, err := ProxyRepo(r).Get(request.DestChain).ExecuteMatch(types.ExecuteMatchParams{
		SrcChain:    *srcChain,
		DestChain:   *destChain,
		Order:       order,
		Match:       match,
		OrderStatus: orderStatus,
		MatchStatus: matchStatus,
		Receiver:    request.Receiver,
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to create match")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ape.Render(w, models.NewTxResponse(tx, *destChain))
}
