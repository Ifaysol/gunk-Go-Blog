package postgres

import (
	"context"
	"grpc-blog/blog/storage"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateCategory(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      storage.Category
		want    int64
		wantErr bool
	}{
		{
			name: "CREATE_CATEGORY_SUCCESS",
			in: storage.Category{
				CategoryName: "This is category",
				CategoryDescription: "This is category description",
			},
			want: 1,
		},
		{
			name: "CREATE_CATEGORY_SUCCESS",
			in: storage.Category{
				CategoryName: "This is category 2",
				CategoryDescription: "This is category description 2",
			},
			want: 2,
		},

		{
			name: "FAILED_DUPLICATE_CATEGORYNAME",
			in: storage.Category{
				CategoryName: "This is category",
				CategoryDescription: "This is category description",
			},
			wantErr: true,
		},

	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Create(context.TODO(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCategory(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      int64
		want    *storage.Category
		wantErr bool
	}{
		{
			name: "GET_CATEGORY_SUCCESS",
			in: 1,
			want:  &storage.Category{
				ID:                  1,
				CategoryName:        "This is category",
				CategoryDescription: "This is category description",
			},
		},
		{
			name: "FAILED_CATEGORY_DOES_NOT_EXIST",
			in: 100,
			wantErr: true,
		},
    }
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
		

			got, err := s.Get(context.TODO(), tt.in) 
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Diff: got -, want + = %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestUpdateCategory(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      storage.Category
		want    *storage.Category
		wantErr bool
	}{
		{
			name: "UPDATE_CATEGORY_SUCCESS",
			in: storage.Category{
				ID: 1,
				CategoryName: "This is category update",
				CategoryDescription: "This is category description update",
			},
			want: &storage.Category{
				ID: 1,
				CategoryName: "This is category update",
				CategoryDescription: "This is category description update",
			},
		},

	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			 err := s.Update(context.TODO(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDeleteCategory(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      int64
		wantErr bool
	}{
		{
			name: "DELETE_CATEGORY_SUCCESS",
			in: 1,
		},
		{
			name: "FAILED_TO_DELETE_CATEGORY_ID_DOES_NOT_EXIST",
			in: 100,
			wantErr: true,
		},
    }
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := s.Delete(context.TODO(), tt.in) 
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
		


