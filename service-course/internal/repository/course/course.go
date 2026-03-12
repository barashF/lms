package course

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/barashF/lms/service-course/internal/logger"
	"github.com/barashF/lms/service-course/internal/model"
)

type Repository struct {
	manager manager
	logger  logger.Logger
}

func NewRepository(m manager, l logger.Logger) *Repository {
	return &Repository{
		manager: m,
		logger:  l,
	}
}

func (r *Repository) Create(ctx context.Context, courseModel model.Course) (uuid.UUID, error) {
	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		r.logger.Error("get conn from transaction manager", logger.NewField("error", err))
		return uuid.Nil, fmt.Errorf("get conn from transaction manager: %w", err)
	}

	query := "INSERT INTO courses (id, title, description, price, type, author) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = conn.Exec(ctx, query,
		courseModel.ID, courseModel.Title, courseModel.Description, courseModel.Price, courseModel.Type, courseModel.Author)
	if err != nil {
		r.logger.Error("database error", logger.NewField("error", err))
		return uuid.Nil, fmt.Errorf("database error: %w", err)
	}

	return courseModel.ID, nil
}

func (r *Repository) FetchByID(ctx context.Context, id uuid.UUID) (*model.Course, error) {
	var course model.Course

	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		r.logger.Error("get conn from transaction manager", logger.NewField("error", err))
		return nil, fmt.Errorf("get conn from transaction manager: %w", err)
	}

	err = conn.QueryRow(ctx, `
		SELECT id, title, description, price, type, author, created_at, updated_at FROM courses
		WHERE id = $1
		`, id).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.Price,
		&course.Type,
		&course.Author,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Error("course not found")
			return nil, model.ErrNoCourseFound
		}
		r.logger.Error("database error", logger.NewField("error", err))
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &course, nil
}

func (r *Repository) FetchAll(ctx context.Context) ([]*model.Course, error) {
	var courses []*model.Course

	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		r.logger.Error("get conn from transaction manager", logger.NewField("error", err))
		return nil, fmt.Errorf("get conn from transaction manager: %w", err)
	}

	rows, err := conn.Query(ctx, `
		SELECT id, title, description, price, type, author, created_at, updated_at
		FROM courses`)
	if err != nil {
		r.logger.Error("database error", logger.NewField("error", err))
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var course model.Course

		err = rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.Price,
			&course.Type,
			&course.Author,
			&course.CreatedAt,
			&course.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("error reading data", logger.NewField("error", err))
			return nil, fmt.Errorf("error reading data: %w", err)
		}

		courses = append(courses, &course)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("database error", logger.NewField("error", err))
		return nil, fmt.Errorf("database error: %w", err)
	}

	if courses == nil {
		courses = []*model.Course{}
	}

	return courses, nil
}
