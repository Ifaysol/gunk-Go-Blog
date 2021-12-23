package post

import (
	"context"
	"grpc-blog/blog/storage"
	ppb "grpc-blog/gunk/v1/post"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CreatePost(ctx context.Context, req *ppb.CreatePostRequest) (*ppb.CreatePostResponse, error) {
	// Need to validate request
	post := storage.Post{
		CatID:     req.Post.CatID,
		PostName:  req.Post.PostName,
		PostImage: req.Post.PostImage,
	}
	id, err := s.pcore.CreatePost(context.Background(), post)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create post.")
	}
	return &ppb.CreatePostResponse{
		ID: id,
	}, nil
}

func (s Svc) ListPost(ctx context.Context, req *ppb.ListPostRequest) (*ppb.ListPostResponse, error) {
	id, err := s.pcore.ListPost(context.Background())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create post: %s", err)
	}
	var pl []*ppb.Post
	for _, value := range id {
		pl = append(pl, &ppb.Post{
			ID:           value.ID,
			CatID:        value.CatID,
			PostName:     value.PostName,
			PostImage:    value.PostImage,
			CategoryName: value.CategoryName,
		})
	}
	return &ppb.ListPostResponse{
		Post: pl,
	}, nil
}

func (s Svc) GetPost(ctx context.Context, req *ppb.GetPostRequest) (*ppb.GetPostResponse, error) {

	pl, err := s.pcore.GetPost(context.Background(), req.GetID())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get post.")
	}

	return &ppb.GetPostResponse{
		Post: &ppb.Post{
			ID:        pl.ID,
			CatID:     pl.CatID,
			PostName:  pl.PostName,
			PostImage: pl.PostImage,
			CategoryName: pl.CategoryName,

		},
	}, nil
}

func (s Svc) UpdatePost(ctx context.Context, req *ppb.UpdatePostRequest) (*ppb.UpdatePostResponse, error) {
	post := storage.Post{
		ID:        req.Post.ID,
		CatID:     req.Post.CatID,
		PostName:  req.Post.PostName,
		PostImage: req.Post.PostImage,
	}
	err := s.pcore.UpdatePost(context.Background(), post)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to Update post.")
	}
	return &ppb.UpdatePostResponse{}, nil
}

func (s Svc) DeletePost(ctx context.Context, req *ppb.DeletePostRequest) (*ppb.DeletePostResponse, error) {
	err := s.pcore.DeletePost(context.Background(), req.GetID())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to Delete post.")
	}
	return &ppb.DeletePostResponse{}, nil
}
