package data

import (
	"encoding/json"

	"github.com/Swapica/swapica-svc/resources"
)

type ChainsQ interface {
	New() ChainsQ
	Select() ([]Chain, error)
	Get() (*Chain, error)

	FilterByID(ids ...string) ChainsQ
	FilterByType(types ...resources.ChainType) ChainsQ
}

type Chain struct {
	ID          string              `fig:"id,required"`
	Name        string              `fig:"name,required"`
	Icon        *string             `fig:"icon"`
	Type        resources.ChainType `fig:"type,required"`
	ChainParams json.RawMessage     `fig:"chain_params"`

	Confirmations   int    `fig:"confirmations,required"`
	SwapContract    string `fig:"swap_contract,required"`
	RelayerContract string `fig:"relayer_contract,required"`
	RpcEndpoint     string `fig:"rpc_endpoint,required"`
	// Relation
	Tokens []TokenChain
}
