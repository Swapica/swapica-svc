package handlers

import (
	"github.com/Swapica/swapica-svc/internal/service/models"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetTokenChainList(w http.ResponseWriter, r *http.Request) {
	tokenChainList, err := TokenChainsQ(r).Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get chain list")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, models.NewTokenChainListResponse(tokenChainList))
}
