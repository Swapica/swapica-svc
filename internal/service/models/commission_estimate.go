package models

import "math/big"

type CommissionEstimateResponse struct {
	LowerCommission    string `json:"lower_commission"`
	LowerCommissionUsd string `json:"lower_commission_usd"`
	UpperCommission    string `json:"upper_commission"`
	UpperCommissionUsd string `json:"upper_commission_usd"`
}

func NewCommissionEstimateResponse(lower, lowerUsd, upper, upperUsd *big.Float) CommissionEstimateResponse {
	response := CommissionEstimateResponse{
		LowerCommission:    lower.String(),
		LowerCommissionUsd: lowerUsd.String(),
		UpperCommission:    upper.String(),
		UpperCommissionUsd: upperUsd.String(),
	}

	return response
}
