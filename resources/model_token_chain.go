/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type TokenChain struct {
	Key
	Attributes    TokenChainAttributes    `json:"attributes"`
	Relationships TokenChainRelationships `json:"relationships"`
}
type TokenChainResponse struct {
	Data     TokenChain `json:"data"`
	Included Included   `json:"included"`
}

type TokenChainListResponse struct {
	Data     []TokenChain `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustTokenChain - returns TokenChain from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustTokenChain(key Key) *TokenChain {
	var tokenChain TokenChain
	if c.tryFindEntry(key, &tokenChain) {
		return &tokenChain
	}
	return nil
}
