package service

import (
	"github.com/Swapica/swapica-svc/internal/converter"
	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/internal/data/mem"
	"github.com/Swapica/swapica-svc/internal/proxy"
	"github.com/Swapica/swapica-svc/internal/service/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router(proxyRepo proxy.ProxyRepo, chains data.ChainsQ, tokens data.TokensQ) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxSigner(s.cfg.Signer()),
			handlers.CtxChainsQ(chains),
			handlers.CtxTokensQ(tokens),
			handlers.CtxTokenChainsQ(mem.NewTokenChainsQ(s.cfg.TokenChains())),
			handlers.CtxProxyRepo(proxyRepo),
			handlers.CtxConverter(converter.NewConverter()),
		),
	)
	r.Route("/v1", func(r chi.Router) {
		r.Route("/create", func(r chi.Router) {
			r.Post("/match", handlers.CreateMatch)
			r.Post("/order", handlers.CreateOrder)
		})
		r.Route("/cancel", func(r chi.Router) {
			r.Post("/match", handlers.CancelMatch)
			r.Post("/order", handlers.CancelOrder)
		})
		r.Route("/execute", func(r chi.Router) {
			r.Post("/match", handlers.ExecuteMatch)
			r.Post("/order", handlers.ExecuteOrder)
		})
		r.Get("/chains", handlers.GetChainList)
		r.Get("/tokens", handlers.GetTokenList)
		r.Post("/approve", handlers.Approve)
		r.Get("/commission_estimate", handlers.GetCommissionEstimate)
	})

	return r
}
