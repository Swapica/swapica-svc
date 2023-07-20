/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type GetCommissionEstimateRequest struct {
	// uint256 value, multiplied by 10^token_decimals
	AmountToBuy string `json:"amount_to_buy"`
	DestChain   string `json:"dest_chain"`
	// address of token contract
	TokenToBuy string `json:"token_to_buy"`
}
