package financialservice

import (
	"fmt"
	"math"

	types "github.com/brianwu291/go-learn/types"
)

type (
	FinancialService interface {
		CalculateFinancial(req types.FinancialRawInfo, roundingDigits int) (types.FinancialResultInfo, error)
	}
 	financialService struct {}
)

func NewFinancialService() FinancialService {
	return &financialService{}
}

func (s *financialService) roundFloat(val float64, digits int) (float64, error) {
	if digits <= 0 {
		return 0, fmt.Errorf("invalid digits: must be greater than 0")
	}

	var digitBase = math.Pow(10, float64(digits))

	roundedRatio := math.Round(val * digitBase) / digitBase

	return roundedRatio, nil
}

func (s *financialService) CalculateFinancial(req types.FinancialRawInfo, roundingDigits int) (types.FinancialResultInfo, error) {
	ebt := float64(req.Revenue - req.Expenses)
	profit := ebt * (1 - req.TaxRate)
	ratio := ebt / profit

	roundedProfit, err := s.roundFloat(profit, roundingDigits)
	if err != nil {
		return types.FinancialResultInfo{}, fmt.Errorf("failed to round profit: %w", err)
	}

	roundedRatio, err := s.roundFloat(ratio, roundingDigits)
	if err != nil {
		return types.FinancialResultInfo{}, fmt.Errorf("failed to round ratio: %w", err)
	}

	return types.FinancialResultInfo{
		Profit: roundedProfit,
		Ratio:  roundedRatio,
	}, nil
}
