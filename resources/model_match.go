/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Match struct {
	Account       common.Address `json:"account"`
	AmountToSell  *big.Int       `json:"amount_to_sell,omitempty"`
	Id            *big.Int       `json:"id,omitempty"`
	OriginChain   *big.Int       `json:"origin_chain,omitempty"`
	OriginOrderId *big.Int       `json:"origin_order_id,omitempty"`
	TokenToSell   common.Address `json:"token_to_sell"`
}
