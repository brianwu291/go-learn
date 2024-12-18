package financialhandler

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	gin "github.com/gin-gonic/gin"

	constants "github.com/brianwu291/go-learn/constants"
	financialservice "github.com/brianwu291/go-learn/services/financial"
	types "github.com/brianwu291/go-learn/types"
)

type (
	FinancialHandler struct {
		service financialservice.FinancialService
	}
)

func NewFinancialHandler(service financialservice.FinancialService) *FinancialHandler {
	return &FinancialHandler{
		service: service,
	}
}

func (h *FinancialHandler) Calculate(c *gin.Context) {
	var request types.FinancialRawInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("Error type: %T\n", err)
		fmt.Printf("Error details: %+v\n", err)
		switch {
		case errors.Is(err, io.EOF):
			c.AbortWithStatusJSON(http.StatusBadRequest, types.BadRequestResponse{Message: "Request body is empty"})
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, types.BadRequestResponse{Message: fmt.Sprintf("Invalid request: %v", err)})
		}
		return
	}

	roundingDigits := 2
	response, err := h.service.CalculateFinancial(request, roundingDigits)
	if err != nil {
		internalServerErr := fmt.Errorf(constants.InternalServerErrorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.InternalServerErrorResponse{Message: internalServerErr.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
