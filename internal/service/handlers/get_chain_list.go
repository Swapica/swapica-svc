package handlers

import (
	"net/http"

	"github.com/Swapica/swapica-svc/internal/service/models"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetChainList(w http.ResponseWriter, r *http.Request) {
	chainList, err := ChainsQ(r).Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get chain list")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, models.NewChainListResponse(chainList))
}
