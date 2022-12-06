package models

import (
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/resources"
)

func newChainKey(id string) resources.Key {
	return resources.Key{
		ID:   id,
		Type: resources.CHAIN,
	}
}

func newChainModel(value data.Chain) *resources.Chain {
	return &resources.Chain{
		Key: newChainKey(value.ID),
		Attributes: resources.ChainAttributes{
			Name:        value.Name,
			Icon:        value.Icon,
			ChainType:   value.Type,
			ChainParams: value.ChainParams,
		},
	}
}

func newChainModelWithRelation(value data.Chain) *resources.Chain {
	model := newChainModel(value)
	model.Relationships = resources.ChainRelationships{
		Tokens: resources.RelationCollection{
			Data: make([]resources.Key, len(value.Tokens)),
		},
	}
	for i, relation := range value.Tokens {
		model.Relationships.Tokens.Data[i] = newTokenKey(relation.TokenID)
	}

	return model
}

func NewChainListResponse(chains []data.Chain, tokens []data.Token) resources.ChainListResponse {
	response := resources.ChainListResponse{
		Data:     make([]resources.Chain, len(chains)),
		Included: resources.Included{},
	}

	for i, chain := range chains {
		response.Data[i] = *newChainModelWithRelation(chain)
	}

	for _, token := range tokens {
		response.Included.Add(newTokenModel(token))
	}

	return response
}
