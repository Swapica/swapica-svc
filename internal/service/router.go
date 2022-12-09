package service

import (
	"github.com/Swapica/swapica-svc/internal/data/mem"
	"github.com/Swapica/swapica-svc/internal/proxy"
	"github.com/Swapica/swapica-svc/internal/service/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	proxyRepo, err := proxy.NewProxyRepo(s.cfg.Chains(), s.cfg.Signer())
	if err != nil {
		panic(err)
	}

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxSigner(s.cfg.Signer()),
			handlers.CtxChainsQ(mem.NewChainsQ(s.cfg.Chains())),
			handlers.CtxTokensQ(mem.NewTokenQ(s.cfg.Tokens())),
			handlers.CtxTokenChainsQ(mem.NewTokenChainsQ(s.cfg.TokenChains())),
			handlers.CtxProxyRepo(proxyRepo),
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
	})

	return r
}
