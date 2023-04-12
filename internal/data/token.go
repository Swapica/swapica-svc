package data

import (
	"github.com/Swapica/swapica-svc/resources"
)

type TokensQ interface {
	New() TokensQ
	Select() ([]Token, error)
	Get() (*Token, error)
	FilterByID(ids ...string) TokensQ
	FilterByType(types ...resources.TokenType) TokensQ
}

type Token struct {
	ID        string              `fig:"id,required"`
	Name      string              `fig:"name,required"`
	Symbol    string              `fig:"symbol,required"`
	Icon      *string             `fig:"icon"`
	Type      resources.TokenType `fig:"type,required"`
	MaxAmount uint64              `fig:"max_amount"`
	// Relation
	TokenChains []TokenChain `fig:"chains,required"`
}
