/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CancelMatchRequest struct {
	// The identifier of the chain in which the match was placed.
	ChainId string `json:"chain_id"`
	// The identifier of the match that was selected by the user to cancel it.
	MatchId int `json:"match_id"`
	// The address of the sender
	Sender string `json:"sender"`
}
