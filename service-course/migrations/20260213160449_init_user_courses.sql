-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_courses (
	id UUID PRIMARY KEY,
	course_id UUID NOT NULL REFERENCES courses(id),
	user_id UUID NOT NULL,
	status TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_courses;
-- +goose StatementEnd
