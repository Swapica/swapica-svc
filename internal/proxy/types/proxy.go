package types

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/amount"
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
)

type Proxy interface {
	CreateOrder(params CreateOrderParams) (interface{}, error)
	CancelOrder(params CancelOrderParams) (interface{}, error)
	ExecuteOrder(params ExecuteOrderParams) (interface{}, error)
	CreateMatch(params CreateMatchParams) (interface{}, error)
	CancelMatch(params CancelMatchParams) (interface{}, error)
	ExecuteMatch(params ExecuteMatchParams) (interface{}, error)

	Approve(params ApproveParams) (interface{}, error)

	GetOrder(id *big.Int) (swapica.SwapicaOrder, error)
	GetMatch(id *big.Int) (swapica.SwapicaMatch, error)
	GetOrderStatus(id *big.Int) (swapica.SwapicaStatus, error)
	GetMatchStatus(id *big.Int) (swapica.SwapicaStatus, error)
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
	Order       swapica.SwapicaOrder
	OrderStatus swapica.SwapicaStatus
}

type ExecuteOrderParams struct {
	Sender      string
	SrcChain    data.Chain
	DestChain   data.Chain
	Order       swapica.SwapicaOrder
	Match       swapica.SwapicaMatch
	OrderStatus swapica.SwapicaStatus
	MatchStatus swapica.SwapicaStatus
	Receiver    string
}

type CreateMatchParams struct {
	SrcChain    data.Chain
	DestChain   data.Chain
	Order       swapica.SwapicaOrder
	OrderStatus swapica.SwapicaStatus
	Sender      string
}

type CancelMatchParams struct {
	Sender      string
	SrcChain    data.Chain
	DestChain   data.Chain
	Order       swapica.SwapicaOrder
	Match       swapica.SwapicaMatch
	OrderStatus swapica.SwapicaStatus
	MatchStatus swapica.SwapicaStatus
}

type ExecuteMatchParams struct {
	Sender      string
	SrcChain    data.Chain
	DestChain   data.Chain
	Order       swapica.SwapicaOrder
	Match       swapica.SwapicaMatch
	OrderStatus swapica.SwapicaStatus
	MatchStatus swapica.SwapicaStatus
	Receiver    string
}

type ApproveParams struct {
	Sender       string
	Chain        data.Chain
	TokenAddress string
	TokenType    string
}
