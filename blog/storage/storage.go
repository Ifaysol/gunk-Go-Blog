package storage

type Category struct {
	ID                  int64  `db:"id"`
	CategoryName        string `db:"categoryname"`
	CategoryDescription string `db:"categorydescription"`
}

type Post struct {
	ID           int64  `db:"id"`
	CatID        int64  `db:"catid"`
	PostName     string `db:"postname"`
	PostImage    string `db:"image"`
	CategoryName string `db:"categoryname"`
}
