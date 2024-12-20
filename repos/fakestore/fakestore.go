package fakestorerepo

import (
	"context"
	"fmt"

	"github.com/brianwu291/go-learn/constants"
	"github.com/brianwu291/go-learn/httpclient"
	"github.com/brianwu291/go-learn/types"
)

type (
	FakeStoreRepo struct {
		client *httpclient.Client
	}
)

const (
	categoryPath                = "/products/categories"
	categoryProductPathTemplate = "/products/category/%s"
	productPathTemplate         = "/products/%d"
)

func NewFakeStoreRepo() *FakeStoreRepo {
	return &FakeStoreRepo{
		client: httpclient.NewClient(
			httpclient.WithBaseURL(constants.FakeStoreBaseURL),
		),
	}
}

func (f *FakeStoreRepo) GetCategories(ctx context.Context) ([]types.Category, error) {
	var categories []types.Category
	err := f.client.Get(ctx, categoryPath, &categories)
	if err != nil {
		return nil, fmt.Errorf("get all categories failed: %w", err)
	}
	return categories, nil
}

func (f *FakeStoreRepo) GetProductsByCategory(ctx context.Context, category types.Category) ([]types.Product, error) {
	var products []types.Product
	categoryProductPath := fmt.Sprintf(categoryProductPathTemplate, category)
	err := f.client.Get(ctx, categoryProductPath, &products)
	if err != nil {
		return nil, fmt.Errorf("get category: %s products failed: %w", category, err)
	}
	return products, nil
}

func (f *FakeStoreRepo) GetProduct(ctx context.Context, id int64) (*types.Product, error) {
	var product types.Product
	productPath := fmt.Sprintf(productPathTemplate, id)
	err := f.client.Get(ctx, productPath, &product)
	if err != nil {
		return nil, fmt.Errorf("get product failed %d: %w", id, err)
	}
	return &product, nil
}
