package enums

type State uint8

const (
	None State = iota
	AwaitingMatch
	AwaitingFinalization
	Canceled
	Executed
)

const (
	TokenTypeNative  = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
	TokenTypeErc20   = "erc20"
	TokenTypeErc721  = "erc721"
	TokenTypeErc1155 = "erc1155"
)
