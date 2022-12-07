/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import (
	"github.com/Swapica/swapica-svc/internal/proxy/evm"
	"math/big"
)

type Status struct {
	ExecutedBy *big.Int  `json:"executed_by,omitempty"`
	State      evm.State `json:"state"`
}
