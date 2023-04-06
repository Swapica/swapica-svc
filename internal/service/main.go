package service

import (
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/internal/data/mem"
	"github.com/Swapica/swapica-svc/internal/proxy"
	"github.com/Swapica/swapica-svc/internal/runner"
	"net"
	"net/http"

	"github.com/Swapica/swapica-svc/internal/config"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log      *logan.Entry
	copus    types.Copus
	listener net.Listener
	cfg      config.Config
}

func (s *service) run(proxyRepo proxy.ProxyRepo, chains data.ChainsQ, tokens data.TokensQ) error {
	s.log.Info("Service started")
	r := s.router(proxyRepo, chains, tokens)

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func newService(cfg config.Config) *service {
	return &service{
		log:      cfg.Log(),
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
		cfg:      cfg,
	}
}

func Run(cfg config.Config) {
	proxyRepo, err := proxy.NewProxyRepo(cfg.Chains(), cfg.Signer())
	if err != nil {
		panic(err)
	}

	chains := mem.NewChainsQ(cfg.Chains())
	tokens := mem.NewTokenQ(cfg.Tokens())

	newRunner, err := runner.NewRunner(cfg, proxyRepo, chains, tokens)
	if err != nil {
		panic(err)
	}

	if cfg.RunnerCfg().SendAutoExecute {
		newRunner.ExecuteOrders()
	}
	go newRunner.ListenNewOrders()

	if err := newService(cfg).run(proxyRepo, chains, tokens); err != nil {
		panic(err)
	}
}
