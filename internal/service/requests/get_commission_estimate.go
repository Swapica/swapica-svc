package requests

import (
	"net/http"
)

type GetCommissionEstimateRequest struct {
	AmountToBuy string
	DestChain   string
	TokenToBuy  string
}

func NewGetCommissionEstimateRequest(r *http.Request) (GetCommissionEstimateRequest, error) {
	request := GetCommissionEstimateRequest{}
	request.AmountToBuy = r.URL.Query().Get("amount_to_buy")
	request.DestChain = r.URL.Query().Get("dest_chain")
	request.TokenToBuy = r.URL.Query().Get("token_to_buy")

	return request, nil
}
