-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    url_shortener (
        short_url VARCHAR(255) PRIMARY KEY,
        long_url VARCHAR(255) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expired_at bigint NOT NULL
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE url_shortener;

-- +goose StatementEnd