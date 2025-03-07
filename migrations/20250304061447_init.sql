-- +goose Up
-- +goose StatementBegin
create table songs (
    id serial,
    song varchar not null ,
    author varchar not null ,
    release_date varchar not null ,
    text varchar,
    link varchar
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE songs
-- +goose StatementEnd
