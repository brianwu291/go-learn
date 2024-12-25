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
	response, err := h.service.GetCategories(c, false)
	if err != nil {
		internalServerErr := fmt.Errorf(constants.InternalServerErrorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.InternalServerErrorResponse{Message: internalServerErr.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *FakeStoreHandler) GetAllCategoriesProducts(c *gin.Context) {
	allCategories, err := h.service.GetCategories(c, false)
	if err != nil {
		fmt.Printf("failed to get all categories: %+v", err)
		internalServerErr := fmt.Errorf(constants.InternalServerErrorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.InternalServerErrorResponse{Message: internalServerErr.Error()})
		return
	}

	maxWorkers := 5
	jobs := make(chan types.Category, len(allCategories))
	allProducts := make(chan []types.Product, len(allCategories))
	errors := make(chan error, 1)
	for i := 0; i < maxWorkers; i += 1 {
		go func() {
			for job := range jobs {
				categoryProducts, err := h.service.GetProductsByCategory(c, job)
				if err != nil {
					errors <- err
					break
				}
				allProducts <- categoryProducts
			}
		}()
	}
	for _, category := range allCategories {
		jobs <- category
	}
	close(jobs)

	var combinedProducts []types.Product
	var combinedErrors []error
	expectedResults := len(allCategories)
	for i := 0; i < expectedResults; i += 1 {
		select {
		case err := <-errors:
			combinedErrors = append(combinedErrors, err)
			errorMessage := fmt.Sprintf("%s: err when getting categories products: %+v",
				constants.InternalServerErrorMessage, combinedErrors)
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				types.InternalServerErrorResponse{Message: errorMessage})
			return
		case products := <-allProducts:
			combinedProducts = append(combinedProducts, products...)
		}
	}

	c.JSON(http.StatusOK, combinedProducts)
}
