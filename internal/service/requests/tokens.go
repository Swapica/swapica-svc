package requests

import (
	"net/http"

	"github.com/Swapica/swapica-svc/resources"
	"gitlab.com/distributed_lab/urlval"
)

type TokensRequest struct {
	FilterType         []resources.TokenType `filter:"token_type"`
	IncludeChains      bool                  `include:"chain"`
	IncludeTokenChains bool                  `include:"token_chain"`
}

func NewTokensRequest(r *http.Request) (TokensRequest, error) {
	request := TokensRequest{}
	err := urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
