package state

type State uint8

const (
	None State = iota
	AwaitingMatch
	AwaitingFinalization
	Canceled
	Executed
)
