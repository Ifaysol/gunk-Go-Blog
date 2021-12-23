package post

import (
	"context"
	"grpc-blog/blog/storage"
	"grpc-blog/blog/storage/postgres"
)

// type postStore interface {
// 	CreatePost(context.Context,storage.Post) (int64, error)
// 	ListPost(context.Context) ([]storage.Post, error)
// 	GetPost(context.Context,int64) (storage.Post, error)
// 	UpdatePost(context.Context, storage.Post) error
// 	DeletePost(context.Context, int64) error
// }

type PostSvc struct {
	store *postgres.Storage
}

func NewPostSvc(s *postgres.Storage) *PostSvc {
	return &PostSvc{
		store: s,
	}
}

func (ps PostSvc) CreatePost(ctx context.Context, p storage.Post) (int64, error) {
	return ps.store.CreatePost(ctx, p)
}

func (ps PostSvc) ListPost(ctx context.Context) ([]storage.Post, error) {
	return ps.store.ListPost(ctx)
}

func (ps PostSvc) GetPost(ctx context.Context, id int64) (storage.Post, error) {
	return ps.store.GetPost(ctx, id)
}

func (ps PostSvc) UpdatePost(ctx context.Context, p storage.Post) error {
	return ps.store.UpdatePost(ctx, p) 
}


func (ps PostSvc) DeletePost(ctx context.Context, id int64) error {
	return ps.store.DeletePost(ctx, id)
}