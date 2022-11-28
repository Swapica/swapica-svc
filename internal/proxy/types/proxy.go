package types

type Proxy interface {
	CreateOrder() (interface{}, error)
	CancelOrder() (interface{}, error)
	MatchOrder() (interface{}, error)
	CancelMatchOrder() (interface{}, error)
	FinalizeOrder() (interface{}, error)
	ExecuteOrder() (interface{}, error)
	//Approve(tokenChain data.TokenChain, approveFrom string) (interface{}, error)
	//LockFungible(params FungibleLockParams) (interface{}, error)
	//LockNonFungible(params NonFungibleLockParams) (interface{}, error)
	//CheckFungibleLockEvent(txHash string, eventIndex int, tokenChain data.TokenChain) (*FungibleLockEvent, error)
	//CheckNonFungibleLockEvent(txHash string, eventIndex int, tokenChain data.TokenChain) (*NonFungibleLockEvent, error)
	//RedeemFungible(params FungibleRedeemParams) (interface{}, error)
	//RedeemNonFungible(params NonFungibleRedeemParams) (interface{}, error)

	// Balance returns balance of the given token for the given address
	// For fungible tokens it returns amount of tokens
	// For non-fungible tokens returns 1 if the token is owned by the account, 0 otherwise
	// nftId should be not nil for non-fungible tokens
	//Balance(tokenChain data.TokenChain, address string, nftId *string) (amount.Amount, error)
	//BridgeBalance(tokenChain data.TokenChain, nftId *string) (amount.Amount, error)
	//GetNftMetadataUri(tokenChain data.TokenChain, nftId string) (string, error)
	//GetNftMetadata(tokenChain data.TokenChain, nftId string) (*NftMetadata, error)
	//GetDecimals(tokenChain data.TokenChain, chainParams *json.RawMessage) (*uint8, error)
	//ValidateAddress(address string) (bool, error)
}
