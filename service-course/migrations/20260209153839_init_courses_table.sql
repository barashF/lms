-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS courses
(
	id UUID PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	price BIGINT,
	type TEXT NOT NULL,
	author TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS courses;
-- +goose StatementEnd
