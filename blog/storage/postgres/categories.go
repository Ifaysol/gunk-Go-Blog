package postgres

import (
	"context"
	"grpc-blog/blog/storage"
)

const insertCategory = `
	INSERT INTO categories(
		categoryname,
		categorydescription
	) VALUES(
		:categoryname,
		:categorydescription
    )RETURNING id;
`

func (s *Storage) Create(ctx context.Context, c storage.Category) (int64, error) {
	stmt, err := s.db.PrepareNamed(insertCategory)
	if err != nil {
		return 0, err
	}
	var id int64 
	if err := stmt.Get(&id, c); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) Get(ctx context.Context, id int64) (*storage.Category, error) {
	var t storage.Category
	if err := s.db.Get(&t, "SELECT * FROM categories WHERE id =$1", id); err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Storage) List(ctx context.Context) ([]storage.Category, error) {
	var list []storage.Category
	if err := s.db.Select(&list, "SELECT * FROM categories"); err != nil {
		return nil, err
	}
	return list, nil
}

const updateCategory = `
	UPDATE categories 
	SET 
		categoryname = :categoryname,
		categorydescription = :categorydescription
	WHERE 
		id = :id
	RETURNING *;
`

func (s *Storage) Update(ctx context.Context, c storage.Category) error {
	res, err := s.db.PrepareNamed(updateCategory)
	if err != nil {
		return  err
	}

    var uc storage.Category
	if err := res.Get(&uc,c); err!=nil{
		return err
	}
	
	return nil
}

func (s *Storage) Delete(ctx context.Context, id int64) error {
	var t storage.Category
	if err := s.db.Get(&t, "DELETE FROM categories WHERE id =$1 RETURNING *", id); err != nil {
		return err
	}
	return nil
}