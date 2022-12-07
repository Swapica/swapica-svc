/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Order struct {
	Account      common.Address `json:"account"`
	AmountToBuy  *big.Int       `json:"amount_to_buy,omitempty"`
	AmountToSell *big.Int       `json:"amount_to_sell,omitempty"`
	DestChain    *big.Int       `json:"dest_chain,omitempty"`
	Id           *big.Int       `json:"id,omitempty"`
	TokenToBuy   common.Address `json:"token_to_buy"`
	TokenToSell  common.Address `json:"token_to_sell"`
}
