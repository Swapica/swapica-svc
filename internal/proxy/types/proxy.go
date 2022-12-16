package types

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/amount"
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/resources"
)

type Proxy interface {
	CreateOrder(params CreateOrderParams) (interface{}, error)
	CancelOrder(params CancelOrderParams) (interface{}, error)
	ExecuteOrder(params ExecuteOrderParams) (interface{}, error)
	CreateMatch(params CreateMatchParams) (interface{}, error)
	CancelMatch(params CancelMatchParams) (interface{}, error)
	ExecuteMatch(params ExecuteMatchParams) (interface{}, error)

	Approve(params ApproveParams) (interface{}, error)

	GetOrder(id *big.Int) (resources.Order, error)
	GetMatch(id *big.Int) (resources.Match, error)
	GetOrderStatus(id *big.Int) (resources.Status, error)
	GetMatchStatus(id *big.Int) (resources.Status, error)
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

type CancelOrderParams struct {
	Sender      string
	SrcChain    data.Chain
	Order       resources.Order
	OrderStatus resources.Status
}

type ExecuteOrderParams struct {
	Sender      string
	SrcChain    data.Chain
	DestChain   data.Chain
	Order       resources.Order
	Match       resources.Match
	OrderStatus resources.Status
	MatchStatus resources.Status
	Receiver    string
}

type CreateMatchParams struct {
	SrcChain    data.Chain
	DestChain   data.Chain
	Order       resources.Order
	OrderStatus resources.Status
	Sender      string
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
	Sender      string
	SrcChain    data.Chain
	DestChain   data.Chain
	Order       resources.Order
	Match       resources.Match
	OrderStatus resources.Status
	MatchStatus resources.Status
	Receiver    string
}

type ApproveParams struct {
	Sender       string
	Chain        data.Chain
	TokenAddress string
	TokenType    string
}
