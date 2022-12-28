package models

import (
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/resources"
)

func newTokenChainModel(value data.TokenChain) *resources.TokenChain {
	return &resources.TokenChain{
		Attributes: resources.TokenChainAttributes{
			ChainId:         value.ChainID,
			ContractAddress: *value.ContractAddress,
			TokenId:         value.TokenID,
			TokenType:       resources.TokenType(value.TokenType),
		},
	}
}
