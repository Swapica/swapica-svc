package handlers

import (
	"net/http"

	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/internal/service/models"
	"github.com/Swapica/swapica-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetTokenList(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewTokensRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("request is incorrect")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tokensQ := TokensQ(r)
	if len(request.FilterType) > 0 {
		tokensQ.FilterByType(request.FilterType...)
	}
	tokens, err := tokensQ.Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get tokens")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	var chains []data.Chain
	if request.IncludeChains {
		chains, err = ChainsQ(r).FilterByID(chainsId(tokens)...).Select()
		if err != nil {
			Log(r).WithError(err).Error("failed to get chains")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	var tokenChains []data.TokenChain
	if request.IncludeTokenChains {
		tokenChains, err = TokenChainsQ(r).FilterByTokenID(tokensId(tokens)...).Select()
		if err != nil {
			Log(r).WithError(err).Error("failed to get token_chains")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	ape.Render(w, models.NewTokenListResponse(tokens, chains, tokenChains))
}

func chainsId(tokens []data.Token) []string {
	result := make(map[string]struct{})
	for _, token := range tokens {
		for _, chain := range token.Chains {
			result[chain.ChainID] = struct{}{}
		}
	}

	keys := make([]string, 0, len(result))
	for key := range result {
		keys = append(keys, key)
	}

	return keys
}

func tokensId(tokens []data.Token) []string {
	keys := make([]string, 0, len(tokens))
	for _, token := range tokens {
		keys = append(keys, token.ID)
	}

	return keys
}
