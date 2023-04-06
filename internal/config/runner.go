package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"time"
)

type RunnerConfig struct {
	SendAutoExecute bool          `fig:"send_auto_execute,required"`
	Url             string        `fig:"url,required"`
	Ws              string        `fig:"ws,required"`
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
