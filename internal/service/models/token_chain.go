package models

import (
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/resources"
)

func newTokenChainKey(id string) resources.Key {
	return resources.Key{
		ID:   id,
		Type: resources.TOKEN_CHAIN,
	}
}

func newTokenChainModel(value data.TokenChain) *resources.TokenChain {
	return &resources.TokenChain{
		Key: newTokenChainKey(value.ID),
		Attributes: resources.TokenChainAttributes{
			ChainId:         value.ChainID,
			ContractAddress: *value.ContractAddress,
			TokenId:         value.TokenID,
			TokenType:       resources.TokenType(value.TokenType),
		},
	}
}

func newTokenChainModelWithRelation(value data.TokenChain) *resources.TokenChain {
	model := newTokenChainModel(value)
	model.Relationships = resources.TokenChainRelationships{
		Chains: resources.RelationCollection{
			Data: make([]resources.Key, len(value.Chains)),
		},
	}
	for i, chain := range value.Chains {
		model.Relationships.Chains.Data[i] = newChainKey(chain.ID)
	}

	return model
}
