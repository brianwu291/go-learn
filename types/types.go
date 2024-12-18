package types

type (
	BadRequestResponse struct {
		Message string `json:"message"`
	}

	InternalServerErrorResponse struct {
		Message string `json:"message"`
	}

	FinancialRawInfo struct {
		Revenue  int     `json:"revenue" binding:"required,gte=0"`
		Expenses int     `json:"expenses" binding:"required,gte=0"`
		TaxRate  float64 `json:"taxRate" binding:"gte=0,lte=1"`
	}

	FinancialResultInfo struct {
		Profit float64 `json:"profit" binding:"required"`
		Ratio  float64 `json:"ratio" binding:"required"`
	}
)
