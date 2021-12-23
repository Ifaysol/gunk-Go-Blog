package categories

import (
	"context"
	"grpc-blog/blog/storage"
	tpb "grpc-blog/gunk/v1/categories"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) Create(ctx context.Context, req *tpb.CreateCategoryRequest) (*tpb.CreateCategoryResponse, error) {
	// Need to validate request
	category := storage.Category{
		ID:                  req.GetCategory().ID,
		CategoryName:        req.GetCategory().CategoryName,
		CategoryDescription: req.GetCategory().CategoryDescription,
	}
	id, err := s.ccore.Create(context.Background(), category)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create category.")
	}
	return &tpb.CreateCategoryResponse{
		ID: id,
	}, nil
}

func (s Svc) List(ctx context.Context, req *tpb.ListCategoryRequest) (*tpb.ListCategoryResponse, error) {
	id, err := s.ccore.List(context.Background())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create category: %s", err)
	}
	var catl []*tpb.Category
	for _, value := range id {
		catl = append(catl, &tpb.Category{
			ID:                  value.ID,
			CategoryName:        value.CategoryName,
			CategoryDescription: value.CategoryDescription,
		})
	}
	return &tpb.ListCategoryResponse{
		Category: catl,
	}, nil
}

func (s Svc) Get(ctx context.Context, req *tpb.GetCategoryRequest) (*tpb.GetCategoryResponse, error) {

	catl, err := s.ccore.Get(context.Background(), req.GetID())
	if err != nil{
		return nil, status.Error(codes.Internal, "Failed to get category.")
	}

	return  &tpb.GetCategoryResponse{
		Category : &tpb.Category{
			ID: catl.ID,
			CategoryName: catl.CategoryName,
			CategoryDescription: catl.CategoryDescription,
		},
	}, nil
}


func (s Svc) Update(ctx context.Context, req *tpb.UpdateCategoryRequest ) (*tpb.UpdateCategoryResponse, error) {
	category := storage.Category{
		ID: req.GetCategory().ID,
		CategoryName: req.GetCategory().CategoryName,
		CategoryDescription: req.GetCategory().CategoryDescription,
	}
	err := s.ccore.Update(context.Background(), category)
	if err != nil{
		return nil, status.Error(codes.Internal, "failed to Update category.")
	}
	return  &tpb.UpdateCategoryResponse{}, nil
}

func (s Svc) Delete(ctx context.Context, req *tpb.DeleteCategoryRequest ) (*tpb.DeleteCategoryResponse, error) {
	err := s.ccore.Delete(context.Background(), req.GetID())
	if err != nil{
		return nil, status.Error(codes.Internal, "Failed to Delete category.")
	}
	return  &tpb.DeleteCategoryResponse{}, nil
}
