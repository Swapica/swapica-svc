package handlers

import (
	"context"
	"net/http"

	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/internal/proxy"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	signerCtxKey
	tokensQCtxKey
	chainsQCtxKey
	tokenChainsQCtxKey
	proxyRepoCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxSigner(entry signature.Signer) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, signerCtxKey, entry)
	}
}

func Signer(r *http.Request) signature.Signer {
	return r.Context().Value(signerCtxKey).(signature.Signer)
}

func CtxTokensQ(entry data.TokensQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, tokensQCtxKey, entry)
	}
}

func TokensQ(r *http.Request) data.TokensQ {
	return r.Context().Value(tokensQCtxKey).(data.TokensQ).New()
}

func CtxChainsQ(entry data.ChainsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, chainsQCtxKey, entry)
	}
}

func ChainsQ(r *http.Request) data.ChainsQ {
	return r.Context().Value(chainsQCtxKey).(data.ChainsQ).New()
}

func CtxTokenChainsQ(entry data.TokenChainsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, tokenChainsQCtxKey, entry)
	}
}

func TokenChainsQ(r *http.Request) data.TokenChainsQ {
	return r.Context().Value(tokenChainsQCtxKey).(data.TokenChainsQ).New()
}

func CtxProxyRepo(entry proxy.ProxyRepo) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, proxyRepoCtxKey, entry)
	}
}

func ProxyRepo(r *http.Request) proxy.ProxyRepo {
	return r.Context().Value(proxyRepoCtxKey).(proxy.ProxyRepo)
}
