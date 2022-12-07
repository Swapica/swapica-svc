package types

import (
	"github.com/Swapica/swapica-svc/internal/amount"
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/resources"
	"math/big"
)

type Proxy interface {
	CreateOrder(params CreateOrderParams) (interface{}, error)
	CancelOrder() (interface{}, error)
	ExecuteOrder() (interface{}, error)
	CreateMatch(params CreateMatchParams) (interface{}, error)
	CancelMatch(params CancelMatchParams) (interface{}, error)
	ExecuteMatch(params ExecuteMatchParams) (interface{}, error)
	GetOrder(id *big.Int) (resources.Order, error)
	GetMatch(id *big.Int) (resources.Match, error)
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

type CreateMatchParams struct {
	SrcChain  data.Chain
	DestChain data.Chain
	Order     resources.Order
	Sender    string
}

type CancelMatchParams struct {
	Sender      string
	SrcChain    data.Chain
	DestChain   data.Chain
	Order       resources.Order
	Match       resources.Match
	OrderStatus resources.Status
	MatchStatus resources.Status
}

type ExecuteMatchParams struct {
	Sender    string
	SrcChain  data.Chain
	DestChain data.Chain
	Order     resources.Order
	Match     resources.Match
	Receiver  string
}
