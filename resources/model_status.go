/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import (
	"math/big"

	"github.com/Swapica/swapica-svc/internal/proxy/evm/state"
)

type Status struct {
	ExecutedBy *big.Int    `json:"executed_by,omitempty"`
	State      state.State `json:"state"`
}
