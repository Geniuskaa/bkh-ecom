-- +goose Up
-- +goose StatementBegin
CREATE TABLE banner_clicks (
    click_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    banner_id BIGINT NOT NULL
);

CREATE INDEX idx_banner_clicks_banner_id_time ON banner_clicks (banner_id, click_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE banner_clicks;
-- +goose StatementEnd
