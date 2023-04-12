package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/url"
	"time"
)

type Aggregator struct {
	Endpoint *url.URL `fig:"endpoint,required"`
	Ws       string   `fig:"ws,required"`
}

type RunnerConfig struct {
	SendAutoExecute bool          `fig:"send_auto_execute,required"`
	RelayerEndpoint *url.URL      `fig:"relayer_endpoint,required"`
	Aggregator      Aggregator    `fig:"aggregator,required"`
	Timeout         time.Duration `fig:"timeout,required"`
}

func (c *config) RunnerCfg() RunnerConfig {
	return c.runnerCfg.Do(func() interface{} {
		var cfg RunnerConfig

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "runner")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out runner config"))
		}

		return cfg
	}).(RunnerConfig)
}
