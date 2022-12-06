/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ExecuteOrderRequest struct {
	// The identifier of the chain in which the order was placed.
	ChainId string `json:"chain_id"`
	// The identifier of the match that was selected by the user.
	MatchId int `json:"match_id"`
	// The identifier of the order that was selected by the user to create the match.
	OrderId int `json:"order_id"`
	// The address of the receiver who will get tokens placed in order.
	Receiver string `json:"receiver"`
}
