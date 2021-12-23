package categories

import (
	"context"
	"grpc-blog/blog/storage"
	tpb "grpc-blog/gunk/v1/categories"
)

type categoryCoreStore interface {
	Create(context.Context, storage.Category) (int64, error)
	List(context.Context) ([]storage.Category, error)
	Get(context.Context, int64) (*storage.Category, error)
	Update(context.Context, storage.Category) error
	Delete(context.Context, int64) error
	
}

type Svc struct {
	tpb.UnimplementedCategoryServiceServer
	ccore categoryCoreStore
}

func NewCategoryServer(ccs categoryCoreStore) *Svc {
	return &Svc{
		ccore: ccs,
	}
}