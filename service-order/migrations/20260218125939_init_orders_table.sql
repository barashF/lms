-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL,
	course_id UUID NOT NULL,
	status TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS outbox (
	id UUID PRIMARY KEY,
	aggregate_type TEXT NOT NULL,
	aggregate_id UUID NOT NULL,
	event_type TEXT NOT NULL,
	payload JSONB NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	processed_at TIMESTAMP,
	processed BOOLEAN NOT NULL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
