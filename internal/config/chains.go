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
}

type chainer struct {
	getter kv.Getter
	once   comfig.Once

	chains []data.Chain
}

func (c *chainer) Chains() []data.Chain {
	c.readConfig()
	return c.chains
}

func (c *chainer) readConfig() {
	c.once.Do(func() interface{} {
		cfg := struct {
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
			cfg.Chains[i] = chain
		}

		c.chains = cfg.Chains

		return nil
	})
}
