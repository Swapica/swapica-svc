package models

import (
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/resources"
)

func newTokenKey(id string) resources.Key {
	return resources.Key{
		ID:   id,
		Type: resources.TOKEN,
	}
}

func newTokenModel(value data.Token) *resources.Token {
	return &resources.Token{
		Key: newTokenKey(value.ID),
		Attributes: resources.TokenAttributes{
			Name:      value.Name,
			Symbol:    value.Symbol,
			Icon:      value.Icon,
			TokenType: value.Type,
			MaxAmount: value.MaxAmount,
		},
	}
}

func newTokenModelWithRelation(value data.Token) *resources.Token {
	model := newTokenModel(value)
	model.Relationships = resources.TokenRelationships{
		Chains: resources.RelationCollection{
			Data: make([]resources.Key, len(value.TokenChains)),
		},
	}
	for i, tokenChain := range value.TokenChains {
		model.Relationships.Chains.Data[i] = newTokenChainKey(tokenChain.ID)
	}

	return model
}

func NewTokenListResponse(tokens []data.Token, chains []data.Chain, tokenChains []data.TokenChain) resources.TokenListResponse {
	response := resources.TokenListResponse{
		Data:     make([]resources.Token, len(tokens)),
		Included: resources.Included{},
	}

	for i, token := range tokens {
		response.Data[i] = *newTokenModelWithRelation(token)
	}

	for _, chain := range chains {
		response.Included.Add(newChainModel(chain))
	}

	for _, tokenChain := range tokenChains {
		response.Included.Add(newTokenChainModelWithRelation(tokenChain))
	}

	return response
}
