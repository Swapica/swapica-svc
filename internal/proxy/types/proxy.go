package types

import (
	"github.com/Swapica/swapica-svc/internal/amount"
	"github.com/Swapica/swapica-svc/internal/data"
)

type Proxy interface {
	CreateOrder(params CreateOrderParams) (interface{}, error)
	CancelOrder() (interface{}, error)
	MatchOrder() (interface{}, error)
	CancelMatchOrder() (interface{}, error)
	FinalizeOrder() (interface{}, error)
	ExecuteOrder() (interface{}, error)
}

type CreateOrderParams struct {
	Sender       string
	SrcChain     data.Chain
	TokenToSell  string
	AmountToSell amount.Amount
	TokenToBuy   string
	AmountToBuy  amount.Amount
	DestChain    data.Chain
}
