package fakestorehandler

import (
	"fmt"
	"net/http"

	gin "github.com/gin-gonic/gin"

	constants "github.com/brianwu291/go-learn/constants"
	fakestoreservice "github.com/brianwu291/go-learn/services/fakestore"
	types "github.com/brianwu291/go-learn/types"
)

type (
	FakeStoreHandler struct {
		service fakestoreservice.FakeStoreService
	}
)

func NewFakeStoreHandler(service fakestoreservice.FakeStoreService) *FakeStoreHandler {
	return &FakeStoreHandler{
		service: service,
	}
}

func (h *FakeStoreHandler) GetAllCategories(c *gin.Context) {
	response, err := h.service.GetCategories(c)
	if err != nil {
		internalServerErr := fmt.Errorf(constants.InternalServerErrorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.InternalServerErrorResponse{Message: internalServerErr.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
