package postgres

import (
	"context"
	"grpc-blog/blog/storage"
)

const insertPost = `
	INSERT INTO posts(
		postname,
		catid,
		image
	) VALUES(
		:postname,
		:catid,
		:image
    )RETURNING id;
`

func (s *Storage) CreatePost(ctx context.Context, p storage.Post) (int64, error) {
	stmt, err := s.db.PrepareNamed(insertPost)
	if err != nil {
		return 0, err
	}
	var id int64
	if err := stmt.Get(&id, p); err != nil {
		return 0, err
	}
	return id, nil
}

func (s Storage) GetPost(ctx context.Context, id int64) (storage.Post, error) {
	var t storage.Post
	if err := s.db.Get(&t, "SELECT posts.* , categories.categoryname FROM posts LEFT JOIN categories ON categories.id = posts.catid WHERE posts.id = $1", id); err != nil {
		return t, err
	}
	return t, nil
}

func (s *Storage) ListPost(ctx context.Context) ([]storage.Post, error) {
	var list []storage.Post
	if err := s.db.Select(&list, "SELECT posts.* , categories.categoryname FROM posts LEFT JOIN categories ON categories.id = posts.catid"); err != nil {
		return nil, err
	}
	return list, nil
}

const updatePost = `
	UPDATE posts 
	SET 
		postname = :postname,
		catid = :catid,
		image = :image
	WHERE 
		id = :id
	RETURNING *;
`

func (s *Storage) UpdatePost(ctx context.Context, p storage.Post) error {
	res, err := s.db.PrepareNamed(updatePost)
	if err != nil {
		return err
	}

	var up storage.Post
	if err := res.Get(&up, p); err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeletePost(ctx context.Context, id int64) error {
	var t storage.Post
	return s.db.Get(&t, "DELETE FROM posts WHERE id =$1 RETURNING id", id)
}

