package enrollment

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/barashF/lms/service-course/internal/model"
)

type Repository struct {
	manager manager
}

func NewRepository(m manager) *Repository {
	return &Repository{manager: m}
}

func (r *Repository) Create(ctx context.Context, enrollment model.Enrollment) (uuid.UUID, error) {
	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get conn from transaction manager: %w", err)
	}

	_, err = conn.Exec(ctx, `
		INSERT INTO user_courses (id, course_id, user_id, status)
		VALUES ($1, $2, $3, $4)`,
		enrollment.ID, enrollment.CourseID, enrollment.UserID, enrollment.Status,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("database error: %w", err)
	}

	return enrollment.ID, nil
}

func (r *Repository) FetchByID(ctx context.Context, id uuid.UUID) (*model.Enrollment, error) {
	var enrollment model.Enrollment
	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("get conn from transaction manager: %w", err)
	}

	err = conn.QueryRow(ctx, `
		SELECT id, course_id, user_id, status, created_at, updated_at
		FROM user_courses
		WHERE id = $1`,
		id).Scan(
		&enrollment.ID,
		&enrollment.CourseID,
		&enrollment.UserID,
		&enrollment.Status,
		&enrollment.CreatedAt,
		&enrollment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNoEnrollmentFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &enrollment, nil
}

func (r *Repository) FetchByCourseIDUserID(ctx context.Context, courseID uuid.UUID, userID uuid.UUID) (*model.Enrollment, error) {
	var enrollment model.Enrollment
	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("get conn from transaction manager: %w", err)
	}

	err = conn.QueryRow(ctx, `
		SELECT id, course_id, user_id, status, created_at, updated_at
		FROM user_courses
		WHERE user_id = $1 AND course_id = $2`,
		userID, courseID).Scan(
		&enrollment.ID,
		&enrollment.CourseID,
		&enrollment.UserID,
		&enrollment.Status,
		&enrollment.CreatedAt,
		&enrollment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNoEnrollmentFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &enrollment, nil
}
