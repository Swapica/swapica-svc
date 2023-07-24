package converter

import (
	"math/big"
	"sync"
	"time"

	lru "github.com/golang/groupcache/lru"
)

const (
	defaultItemLifespan = 5 * time.Second
	defaultCacheSize    = 100
)

type priceCache struct {
	itemLifespan time.Duration
	m            sync.RWMutex
	store        *lru.Cache
}

type priceCacheValue struct {
	price          *big.Float
	expirationTime time.Time
}

func (p *priceCache) Add(key string, price *big.Float) {
	p.m.Lock()
	p.store.Add(key, priceCacheValue{
		price:          price,
		expirationTime: time.Now().Add(p.itemLifespan),
	})
	p.m.Unlock()
}

func (p *priceCache) Get(key string) (interface{}, bool) {
	p.m.RLock()
	raw, isExist := p.store.Get(key)
	p.m.RUnlock()
	return raw, isExist
}
