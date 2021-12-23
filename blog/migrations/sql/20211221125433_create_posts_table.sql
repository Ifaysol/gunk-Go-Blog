-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts(
    id SERIAL NOT NULL,
    catid INTEGER NOT NULL,
    postname TEXT NOT NULL UNIQUE,
    image TEXT NOT NULL,


    PRIMARY KEY(id),
    CONSTRAINT fk_category
    FOREIGN KEY(catid)
    REFERENCES categories(id)
    ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF EXISTS posts;
-- +goose StatementEnd
