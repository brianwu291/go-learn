package fakestoreservice

import (
	"context"

	fakestorerepo "github.com/brianwu291/go-learn/repos/fakestore"
	types "github.com/brianwu291/go-learn/types"
)

type (
	FakeStoreService interface {
		GetCategories(ctx context.Context) ([]types.Category, error)
		GetProductsByCategory(ctx context.Context, category types.Category) ([]types.Product, error)
		GetProduct(ctx context.Context, id int64) (*types.Product, error)
	}

	fakeStoreService struct {
		repo *fakestorerepo.FakeStoreRepo
	}
)

func NewFakeStoreService(repo *fakestorerepo.FakeStoreRepo) *fakeStoreService {
	return &fakeStoreService{
		repo: repo,
	}
}

func (s *fakeStoreService) GetCategories(ctx context.Context) ([]types.Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *fakeStoreService) GetProductsByCategory(ctx context.Context, category types.Category) ([]types.Product, error) {
	return s.repo.GetProductsByCategory(ctx, category)
}

func (s *fakeStoreService) GetProduct(ctx context.Context, id int64) (*types.Product, error) {
	return s.repo.GetProduct(ctx, id)
}
