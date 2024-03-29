package handlers

import (
	"math/big"
	"net/http"

	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/Swapica/swapica-svc/internal/service/models"
	"github.com/Swapica/swapica-svc/internal/service/requests"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func CancelMatch(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCancelMatchRequest(r)
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

	match, err := ProxyRepo(r).Get(destChain.ID).GetMatch(big.NewInt(int64(request.MatchId)))
	if err != nil {
		Log(r).WithError(err).Error("failed to get match")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	order, err := ProxyRepo(r).Get(srcChain.ID).GetOrder(match.OriginOrderId)
	if err != nil {
		Log(r).WithError(err).Error("failed to get order")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	params := types.CancelMatchParams{
		Sender:    request.Sender,
		SrcChain:  *srcChain,
		DestChain: *destChain,
		Match:     match,
		Order:     order,
	}

	if request.RawTxData != nil {
		rawTxData, err := hexutil.Decode(*request.RawTxData)
		if err != nil {
			return
		}

		params.RawTxData = &rawTxData
	}

	tx, err := ProxyRepo(r).Get(request.DestChain).CancelMatch(params)
	if err != nil {
		Log(r).WithError(err).Error("failed to create cancel match transaction")
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
