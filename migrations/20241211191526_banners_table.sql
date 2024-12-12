-- +goose Up
-- +goose StatementBegin
CREATE TABLE banners (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE banners;
-- +goose StatementEnd
