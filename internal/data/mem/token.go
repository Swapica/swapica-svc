package mem

import (
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/resources"
)

func NewTokenQ(tokens []data.Token) data.TokensQ {
	return &tokensQ{
		tokens:  tokens,
		filters: make([]tokenFilter, 0),
	}
}

type tokensQ struct {
	tokens  []data.Token
	filters []tokenFilter
}

type tokenFilter func(token data.Token) bool

func (q *tokensQ) New() data.TokensQ {
	return NewTokenQ(q.tokens)
}

func (q *tokensQ) Get() (*data.Token, error) {
	for _, token := range q.tokens {
		if q.filter(token) {
			return &token, nil
		}
	}

	return nil, nil
}

func (q *tokensQ) Select() ([]data.Token, error) {
	result := make([]data.Token, 0, len(q.tokens))
	for _, token := range q.tokens {
		if q.filter(token) {
			result = append(result, token)
		}
	}

	return result, nil
}

func (q *tokensQ) FilterByID(ids ...string) data.TokensQ {
	q.filters = append(q.filters, func(token data.Token) bool {
		return contains(ids, token.ID)
	})
	return q
}

func (q *tokensQ) FilterByType(types ...resources.TokenType) data.TokensQ {
	q.filters = append(q.filters, func(token data.Token) bool {
		return contains(types, token.Type)
	})
	return q
}

func (q *tokensQ) filter(token data.Token) bool {
	for _, filter := range q.filters {
		if !filter(token) {
			return false
		}
	}

	return true
}

func (q *tokensQ) PageTokens(token data.Token) data.TokensQ {
	result := make([]data.Token, 0, len(q.tokens))
	for _, token := range q.tokens {
		if q.filter(token) {
			result = append(result, token)
		}
	}
	return q
}
