package mem

import (
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/resources"
)

func NewChainsQ(chains []data.Chain) data.ChainsQ {
	return &chainsQ{
		chains:  chains,
		filters: make([]chainFilter, 0),
	}
}

type chainsQ struct {
	chains  []data.Chain
	filters []chainFilter
}

type chainFilter func(value data.Chain) bool

func (q *chainsQ) New() data.ChainsQ {
	return NewChainsQ(q.chains)
}

func (q *chainsQ) Get() (*data.Chain, error) {
	for _, value := range q.chains {
		if q.filter(value) {
			return &value, nil
		}
	}

	return nil, nil
}

func (q *chainsQ) Select() ([]data.Chain, error) {
	result := make([]data.Chain, 0, len(q.chains))
	for _, value := range q.chains {
		if q.filter(value) {
			result = append(result, value)
		}
	}

	return result, nil
}

func (q *chainsQ) FilterByID(ids ...string) data.ChainsQ {
	q.filters = append(q.filters, func(value data.Chain) bool {
		return contains(ids, value.ID)
	})
	return q
}

func (q *chainsQ) FilterByType(types ...resources.ChainType) data.ChainsQ {
	q.filters = append(q.filters, func(value data.Chain) bool {
		return contains(types, value.Type)
	})
	return q
}

func (q *chainsQ) filter(value data.Chain) bool {
	for _, filter := range q.filters {
		if !filter(value) {
			return false
		}
	}

	return true
}
