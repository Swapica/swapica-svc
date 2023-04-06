/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "math/big"

type CreateOrderRequest struct {
	// uint256 value, multiplied by 10^token_decimals
	AmountToBuy big.Int `json:"amount_to_buy"`
	// uint256 value, multiplied by 10^token_decimals
	AmountToSell big.Int `json:"amount_to_sell"`
	DestChain    string  `json:"dest_chain"`
	Sender       string  `json:"sender"`
	SrcChain     string  `json:"src_chain"`
	// address of token contract
	TokenToBuy string `json:"token_to_buy"`
	// address of token contract
	TokenToSell string `json:"token_to_sell"`
}
