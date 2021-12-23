package categories

import (
	"context"
	"grpc-blog/blog/storage"
)

type categoryStore interface {
	Create(context.Context,storage.Category) (int64, error)
	List(context.Context) ([]storage.Category, error)
	Get(context.Context,int64) (*storage.Category, error)
	Update(context.Context, storage.Category) error
	Delete(context.Context, int64) error
}

type CategorySvc struct {
	store categoryStore
}

func NewCategorySvc(s categoryStore) *CategorySvc {
	return &CategorySvc{
		store: s,
	}
}

func (cs CategorySvc) Create(ctx context.Context, c storage.Category) (int64, error) {
	return cs.store.Create(ctx, c)
}

func (cs CategorySvc) List(ctx context.Context) ([]storage.Category, error) {
	return cs.store.List(ctx)
}

func (cs CategorySvc) Get(ctx context.Context, id int64) (*storage.Category, error) {
	return cs.store.Get(ctx, id)
}

func (cs CategorySvc) Update(ctx context.Context, c storage.Category) error {
	return cs.store.Update(ctx, c)
}

func (cs CategorySvc) Delete(ctx context.Context, id int64) error {
	return cs.store.Delete(ctx, id)
}