package data

type TokenChainsQ interface {
	New() TokenChainsQ
	Select() ([]TokenChain, error)
	Get() (*TokenChain, error)

	FilterByTokenID(ids ...string) TokenChainsQ
	FilterByChainID(ids ...string) TokenChainsQ
	FilterByBridgingType(types ...BridgingType) TokenChainsQ
}

type TokenChain struct {
	TokenID         string
	ChainID         string       `fig:"chain_id,required"`
	ContractAddress *string      `fig:"contract_address"`
	TokenType       string       `fig:"token_type,required"`
	BridgingType    BridgingType `fig:"bridging_type,required"`
	AutoSend        bool         `fig:"auto_send"`
}

type BridgingType string

const (
	BridgingTypeLP      BridgingType = "liquidity_pool"
	BridgingTypeWrapped BridgingType = "wrapped"
)
