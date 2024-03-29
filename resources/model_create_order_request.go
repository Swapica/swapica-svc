/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateOrderRequest struct {
	// uint256 value, multiplied by 10^token_decimals
	AmountToBuy string `json:"amount_to_buy"`
	// uint256 value, multiplied by 10^token_decimals
	AmountToSell string `json:"amount_to_sell"`
	DestChain    string `json:"dest_chain"`
	Sender       string `json:"sender"`
	SrcChain     string `json:"src_chain"`
	// address of token contract
	TokenToBuy string `json:"token_to_buy"`
	// address of token contract
	TokenToSell string `json:"token_to_sell"`
	UseRelayer  bool   `json:"use_relayer"`
}
