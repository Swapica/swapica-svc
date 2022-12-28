package mem

import (
	"github.com/Swapica/swapica-svc/internal/data"
)

func NewTokenChainsQ(tokenChains []data.TokenChain) data.TokenChainsQ {
	return &tokenChainsQ{
		tokenChains: tokenChains,
		filters:     make([]tokenChainFilter, 0),
	}
}

type tokenChainsQ struct {
	tokenChains []data.TokenChain
	filters     []tokenChainFilter
}

type tokenChainFilter func(value data.TokenChain) bool

func (q *tokenChainsQ) New() data.TokenChainsQ {
	return NewTokenChainsQ(q.tokenChains)
}

func (q *tokenChainsQ) Get() (*data.TokenChain, error) {
	for _, value := range q.tokenChains {
		if q.filter(value) {
			return &value, nil
		}
	}

	return nil, nil
}

func (q *tokenChainsQ) Select() ([]data.TokenChain, error) {
	result := make([]data.TokenChain, 0, len(q.tokenChains))
	for _, value := range q.tokenChains {
		if q.filter(value) {
			result = append(result, value)
		}
	}

	return result, nil
}

func (q *tokenChainsQ) FilterByTokenID(ids ...string) data.TokenChainsQ {
	q.filters = append(q.filters, func(value data.TokenChain) bool {
		return contains(ids, value.TokenID)
	})
	return q
}

func (q *tokenChainsQ) FilterByChainID(ids ...string) data.TokenChainsQ {
	q.filters = append(q.filters, func(value data.TokenChain) bool {
		return contains(ids, value.ChainID)
	})
	return q
}

func (q *tokenChainsQ) filter(value data.TokenChain) bool {
	for _, filter := range q.filters {
		if !filter(value) {
			return false
		}
	}

	return true
}
