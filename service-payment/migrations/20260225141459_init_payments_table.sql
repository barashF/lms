-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payments (
	id UUID PRIMARY KEY,
	order_id UUID NOT NULL,
	user_id UUID NOT NULL,
	course_id UUID NOT NULL,
	value NUMERIC(10, 2) NOT NULL,
	currency TEXT NOT NULL,
	status TEXT NOT NULL,
	payment_id UUID NOT NULL UNIQUE,
	confirmation_token TEXT NOT NULL,
	idempotency_key TEXT NOT NULL UNIQUE,
	paid BOOLEAN DEFAULT FALSE,

)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
