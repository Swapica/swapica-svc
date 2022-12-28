package config

import (
	"github.com/Swapica/swapica-svc/internal/data"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func NewChainer(getter kv.Getter) Chainer {
	return &chainer{
		getter: getter,
	}
}

type Chainer interface {
	Chains() []data.Chain
	Tokens() []data.Token
	TokenChains() []data.TokenChain
}

type chainer struct {
	getter kv.Getter
	once   comfig.Once

	chains      []data.Chain
	tokens      []data.Token
	tokenChains []data.TokenChain
}

func (c *chainer) Chains() []data.Chain {
	c.readConfig()
	return c.chains
}

func (c *chainer) Tokens() []data.Token {
	c.readConfig()
	return c.tokens
}

func (c *chainer) TokenChains() []data.TokenChain {
	c.readConfig()
	return c.tokenChains
}

func (c *chainer) readConfig() {
	c.once.Do(func() interface{} {
		cfg := struct {
			Tokens []data.Token `fig:"tokens,required"`
			Chains []data.Chain `fig:"chains,required"`
		}{}
		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "data")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out tokens and chains"))
		}

		for i, chain := range cfg.Chains {
			chain.Tokens = make([]data.TokenChain, 0)
			cfg.Chains[i] = chain
		}

		tokenChains := make([]data.TokenChain, 0)
		for _, token := range cfg.Tokens {
			for i, tokenChain := range token.Chains {
				tokenChain.TokenID = token.ID
				token.Chains[i] = tokenChain
				tokenChains = append(tokenChains, tokenChain)
				for k, chain := range cfg.Chains {
					if chain.ID == tokenChain.ChainID {
						chain.Tokens = append(chain.Tokens, tokenChain)
						cfg.Chains[k] = chain
					}
				}
			}
		}

		c.tokens = cfg.Tokens
		c.chains = cfg.Chains
		c.tokenChains = tokenChains

		return nil
	})
}
