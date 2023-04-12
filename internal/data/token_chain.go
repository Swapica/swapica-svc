package data

type TokenChainsQ interface {
	New() TokenChainsQ
	Select() ([]TokenChain, error)
	Get() (*TokenChain, error)

	FilterByTokenID(ids ...string) TokenChainsQ
	FilterByChainID(ids ...string) TokenChainsQ
}

type TokenChain struct {
	ID              string
	TokenID         string
	ChainID         string
	ContractAddress *string
	TokenType       string
	MaxAmount       uint64
	// Relation
	Chains []Chain
}
