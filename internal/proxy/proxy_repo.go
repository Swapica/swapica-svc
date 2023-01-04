package proxy

import (
	"fmt"

	"github.com/Swapica/swapica-svc/internal/data"
	"github.com/Swapica/swapica-svc/internal/proxy/evm"
	"github.com/Swapica/swapica-svc/internal/proxy/evm/signature"
	"github.com/Swapica/swapica-svc/internal/proxy/types"
	"github.com/Swapica/swapica-svc/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type ProxyRepo interface {
	Get(chainID string) types.Proxy
}

func NewProxyRepo(chains []data.Chain, signer signature.Signer) (ProxyRepo, error) {
	repo := proxyRepo{
		proxies: make(map[string]types.Proxy),
	}

	for _, c := range chains {
		switch c.Type {
		case resources.EVM:
			proxy, err := evm.NewProxy(c.RpcEndpoint, signer, c.SwapContract, c.Confirmations)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create evm proxy")
			}
			repo.proxies[c.ID] = proxy
		default:
			return nil, errors.New(fmt.Sprintf("Unsupported network type"))
		}
	}

	return &repo, nil
}

type proxyRepo struct {
	proxies map[string]types.Proxy
}

func (r *proxyRepo) Get(chainID string) types.Proxy {
	return r.proxies[chainID]
}
