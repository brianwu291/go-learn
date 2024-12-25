package fakestoreservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brianwu291/go-learn/cache"
	fakestorerepo "github.com/brianwu291/go-learn/repos/fakestore"
	types "github.com/brianwu291/go-learn/types"
)

type (
	FakeStoreService interface {
		GetCategories(ctx context.Context, skipCache bool) ([]types.Category, error)
		GetProductsByCategory(ctx context.Context, category types.Category) ([]types.Product, error)
		GetProduct(ctx context.Context, id int64) (*types.Product, error)
	}

	fakeStoreService struct {
		cacheClient cache.Client
		repo        *fakestorerepo.FakeStoreRepo
	}
)

func NewFakeStoreService(cacheClient cache.Client, repo *fakestorerepo.FakeStoreRepo) *fakeStoreService {
	return &fakeStoreService{
		cacheClient: cacheClient,
		repo:        repo,
	}
}

func (s *fakeStoreService) GetCategories(ctx context.Context, skipCache bool) ([]types.Category, error) {
	if skipCache {
		return s.repo.GetCategories(ctx)
	}

	const categoriesCacheKey = "fakeStore:categories:all"

	if categories, err := s.getCachedCategories(ctx, categoriesCacheKey); err == nil {
		return categories, nil
	}

	categories, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	s.cacheCategories(ctx, categoriesCacheKey, categories)

	return categories, nil
}

func (s *fakeStoreService) getCachedCategories(ctx context.Context, key string) ([]types.Category, error) {
	cachedData, err := s.cacheClient.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var categories []types.Category
	if err := json.Unmarshal([]byte(cachedData), &categories); err != nil {
		fmt.Printf("failed to unmarshal cached categories: %+v", err)
		return nil, err
	}

	return categories, nil
}

func (s *fakeStoreService) cacheCategories(ctx context.Context, key string, categories []types.Category) {
	categoriesJson, err := json.Marshal(categories)
	if err != nil {
		fmt.Printf("failed to marshal categories for cache: %+v", err)
		return
	}

	if err := s.cacheClient.Set(ctx, key, categoriesJson, time.Hour); err != nil {
		fmt.Printf("failed to cache categories: %+v", err)
	}
}

func (s *fakeStoreService) GetProductsByCategory(ctx context.Context, category types.Category) ([]types.Product, error) {
	return s.repo.GetProductsByCategory(ctx, category)
}

func (s *fakeStoreService) GetProduct(ctx context.Context, id int64) (*types.Product, error) {
	return s.repo.GetProduct(ctx, id)
}
