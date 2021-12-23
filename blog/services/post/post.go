package post

import (
	"context"
	"grpc-blog/blog/storage"
	ppb "grpc-blog/gunk/v1/post"
)

type postCoreStore interface {
	CreatePost(context.Context, storage.Post) (int64, error)
	ListPost(context.Context) ([]storage.Post, error)
	GetPost(context.Context, int64) (storage.Post, error)
	UpdatePost(context.Context, storage.Post) error
	DeletePost(context.Context, int64) error
	
}

type Svc struct {
	ppb.UnimplementedPostServiceServer
	pcore postCoreStore
}

func NewPostServer(pcs postCoreStore) *Svc {
	return &Svc{
		pcore: pcs,
	}
}