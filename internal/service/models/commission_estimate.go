package models

import "math/big"

type CommissionEstimateResponse struct {
	LowerCommission string `json:"lower_commission"`
	UpperCommission string `json:"upper_commission"`
}

func NewCommissionEstimateResponse(lower, upper *big.Float) CommissionEstimateResponse {
	response := CommissionEstimateResponse{
		LowerCommission: lower.String(),
		UpperCommission: upper.String(),
	}

	return response
}
