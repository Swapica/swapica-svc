package requests

import (
	"net/http"

	"github.com/go-chi/chi"
)

type GetCommissionEstimateRequest struct {
	AmountToBuy string
	DestChain   string
	TokenToBuy  string
}

func NewGetCommissionEstimateRequest(r *http.Request) (GetCommissionEstimateRequest, error) {
	amountToBuy := chi.URLParam(r, "amount_to_buy")
	destChain := chi.URLParam(r, "dest_chain")
	tokenToBuy := chi.URLParam(r, "token_to_buy")

	request := GetCommissionEstimateRequest{
		AmountToBuy: amountToBuy,
		DestChain:   destChain,
		TokenToBuy:  tokenToBuy,
	}

	return request, nil
}
