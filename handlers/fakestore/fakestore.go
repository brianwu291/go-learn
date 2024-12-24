package fakestorehandler

import (
	"fmt"
	"net/http"
	"sync"

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
		internalServerErr := fmt.Errorf(constants.InternalServerErrorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.InternalServerErrorResponse{Message: internalServerErr.Error()})
		return
	}

	var wg sync.WaitGroup
	allProducts := make(chan []types.Product, len(allCategories))
	errors := make(chan error, len(allCategories))
	for i := 0; i < len(allCategories); i += 1 {
		wg.Add(1)
		go func(category types.Category) {
			defer wg.Done()
			categoryProducts, err := h.service.GetProductsByCategory(c, category)
			if err != nil {
				errors <- err
				return
			}
			allProducts <- categoryProducts
		}(allCategories[i])
	}
	go func() {
		wg.Wait()
		close(allProducts)
		close(errors)
	}()

	var combinedProducts []types.Product
	var combinedErrors []error
	productsOpen, errorsOpen := true, true

	for productsOpen || errorsOpen {
		select {
		case products, ok := <-allProducts:
			if !ok {
				productsOpen = false
				continue
			}
			combinedProducts = append(combinedProducts, products...)

		case err, ok := <-errors:
			if !ok {
				errorsOpen = false
				continue
			}
			combinedErrors = append(combinedErrors, err)
		}
	}

	if len(combinedErrors) > 0 {
		errorMessage := fmt.Sprintf("%s: err when getting categories products: %+v",
			constants.InternalServerErrorMessage, combinedErrors)
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			types.InternalServerErrorResponse{Message: errorMessage})
		return
	}

	c.JSON(http.StatusOK, combinedProducts)
}
