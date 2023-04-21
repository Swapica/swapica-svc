package types

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/generated/swapica"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var ErrWrongSignedTx = errors.New("signed tx does not match tx log")

type Proxy interface {
	CreateOrder(params CreateOrderParams) (interface{}, error)
	CancelOrder(params CancelOrderParams) (interface{}, error)
	ExecuteOrder(params ExecuteOrderParams) (interface{}, error)
	CreateMatch(params CreateMatchParams) (interface{}, error)
	CancelMatch(params CancelMatchParams) (interface{}, error)
	ExecuteMatch(params ExecuteMatchParams) (interface{}, error)

	Approve(params ApproveParams) (interface{}, error)

	GetOrder(id *big.Int) (swapica.ISwapicaOrder, error)
	GetMatch(id *big.Int) (swapica.ISwapicaMatch, error)
}

type CreateOrderParams struct {
	Sender       string
	SrcChain     data.Chain
	TokenToSell  string
	AmountToSell *big.Int
	TokenToBuy   string
	AmountToBuy  *big.Int
	DestChain    data.Chain
	UseRalyer    bool
}

type CancelOrderParams struct {
	Sender   string
	SrcChain data.Chain
	Order    swapica.ISwapicaOrder
}

type ExecuteOrderParams struct {
	Sender    string
	SrcChain  data.Chain
	DestChain data.Chain
	Order     swapica.ISwapicaOrder
	Match     swapica.ISwapicaMatch
	Receiver  string
	RawTxData *[]byte
}

type CreateMatchParams struct {
	SrcChain   data.Chain
	DestChain  data.Chain
	Order      swapica.ISwapicaOrder
	Sender     string
	RawTxData  *[]byte
	UseRelayer bool
}

type CancelMatchParams struct {
	Sender    string
	SrcChain  data.Chain
	DestChain data.Chain
	Order     swapica.ISwapicaOrder
	Match     swapica.ISwapicaMatch
	RawTxData *[]byte
}

type ExecuteMatchParams struct {
	Sender    string
	SrcChain  data.Chain
	DestChain data.Chain
	Order     swapica.ISwapicaOrder
	Match     swapica.ISwapicaMatch
	Receiver  string
	RawTxData *[]byte
}

type ApproveParams struct {
	Sender       string
	Chain        data.Chain
	TokenAddress string
	TokenType    string
}
