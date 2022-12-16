package handlers

import (
	"net/http"

	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/Swapica/swapica-svc/internal/service/models"
	"github.com/Swapica/swapica-svc/internal/service/requests"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func Approve(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewApproveRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("failed to decode request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	chain, err := ChainsQ(r).FilterByID(request.ChainId).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get chain")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if chain == nil {
		Log(r).WithError(err).Debug("chain not found")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"data": errors.New("token that you have sent does not connected to this network or does not exist"),
		})...)
		return
	}

	tx, err := ProxyRepo(r).Get(chain.ID).Approve(types.ApproveParams{
		Sender:       request.Sender,
		Chain:        *chain,
		TokenAddress: request.TokenAddress,
		TokenType:    request.TokenType,
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to approve")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if tx == nil {
		// No approve is needed
		w.WriteHeader(http.StatusNoContent)
		return
	}

	ape.Render(w, models.NewTxResponse(tx, *chain))
}
