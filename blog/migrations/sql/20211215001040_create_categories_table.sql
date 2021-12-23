-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories(
    id SERIAL NOT NULL,
    categoryname TEXT NOT NULL UNIQUE,
    categorydescription TEXT NOT NULL,


    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF EXISTS categories;
-- +goose StatementEnd
