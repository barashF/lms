package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/barashF/lms/service-order/internal/model"
)

type Repository struct {
	manager manager
}

func NewRepository(m manager) *Repository {
	return &Repository{manager: m}
}

func (r *Repository) Create(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get conn from transaction manager: %w", err)
	}

	_, err = conn.Exec(ctx, `
		INSERT INTO orders (id, user_id, course_id, status) VALUES ($1, $2, $3, $4)`,
		order.ID, order.UserID, order.CourseID, order.Status)
	if err != nil {
		return uuid.Nil, fmt.Errorf("database error: %w", err)
	}

	return order.ID, nil
}

func (r *Repository) FetchByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	var order model.Order

	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("get conn from transaction manager: %w", err)
	}

	err = conn.QueryRow(ctx,
		`SELECT id, user_id, course_id, status, created_at, updated_at
		 FROM orders
		 WHERE id = $1`, id,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.CourseID,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNoOrderFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &order, nil
}

func (r *Repository) FetchAll(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order

	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("get conn from transaction manager: %w", err)
	}

	rows, err := conn.Query(ctx, `
		SELECT id, user_id, course_id, status, created_at, updated_at
		 FROM orders`)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.CourseID,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error reading data: %w", err)
		}

		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if orders == nil {
		orders = []model.Order{}
	}
	return orders, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, ordersIDs []uuid.UUID, status string) error {
	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("get conn from transaction manager: %w", err)
	}

	_, err = conn.Exec(ctx, `
		UPDATE orders SET(status, updated_at) = ($1, NOW())
		WHERE id = ANY($2))`,
		status, ordersIDs,
	)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}
