package types

import (
	"github.com/Swapica/swapica-svc/internal/amount"
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/internal/proxy"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Proxy interface {
	CreateOrder(params CreateOrderParams) (interface{}, error)
	CancelOrder() (interface{}, error)
	ExecuteOrder() (interface{}, error)
	CreateMatch() (interface{}, error)
	CancelMatch(params CancelMatchParams, repo proxy.ProxyRepo) (interface{}, error)
	ExecuteMatch() (interface{}, error)

	GetMatch(matchId *big.Int) (struct {
		Id            *big.Int
		OriginOrderId *big.Int
		Account       common.Address
		TokenToSell   common.Address
		AmountToSell  *big.Int
		OriginChain   *big.Int
	}, error)
	GetOrder(orderId *big.Int) (struct {
		Id           *big.Int
		Account      common.Address
		TokenToSell  common.Address
		TokenToBuy   common.Address
		AmountToSell *big.Int
		AmountToBuy  *big.Int
		DestChain    *big.Int
	}, error)
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

type CancelMatchParams struct {
	Sender     string
	SrcChain   data.Chain
	OrgChain   data.Chain
	MatchId    int
	OrderData  []byte
	Signatures [][]byte
}
