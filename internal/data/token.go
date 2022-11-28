package data

type TokensQ interface {
	New() TokensQ
	Select() ([]Token, error)
	Get() (*Token, error)

	FilterByID(ids ...string) TokensQ
}

type Token struct {
	ID     string  `fig:"id,required"`
	Name   string  `fig:"name,required"`
	Symbol string  `fig:"symbol,required"`
	Icon   *string `fig:"icon"`
	// Relation
	Chains []TokenChain `fig:"chains,required"`
}
