package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/barashF/lms/service-order/internal/model"
)

type Repository struct {
	manager manager
}

func NewRepository(m manager) *Repository { return &Repository{manager: m} }

func (r *Repository) Create(ctx context.Context, outbox *model.OutBoxMessage) (uuid.UUID, error) {
	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get conn from transaction manager: %w", err)
	}

	_, err = conn.Exec(ctx, `
		INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload)
		VALUES ($1, $2, $3, $4, $5)
		`, outbox.ID, outbox.AggregateType, outbox.AggregateID, outbox.EventType, outbox.Payload)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error database: %w", err)
	}

	return outbox.ID, nil
}

func (r *Repository) FetchUnprocessed(ctx context.Context, limit int64) ([]*model.OutBoxMessage, error) {
	var messages []*model.OutBoxMessage

	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("get conn from transaction manager: %w", err)
	}

	rows, err := conn.Query(ctx, `
		SELECT id, aggregate_type, aggregate_id, event_type, payload
		FROM outbox
		WHERE processed = false
		ORDER BY created_at ASC
		LIMIT $1
		`, limit)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var message model.OutBoxMessage
		err = rows.Scan(
			&message.ID,
			&message.AggregateType,
			&message.AggregateID,
			&message.EventType,
			&message.Payload,
		)
		if err != nil {
			return nil, fmt.Errorf("error reading data: %w", err)
		}

		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if messages == nil {
		messages = []*model.OutBoxMessage{}
	}

	return messages, nil
}

func (r *Repository) MarkAsProcessed(ctx context.Context, id uuid.UUID) error {
	conn, err := r.manager.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("get conn from transaction manager: %w", err)
	}

	result, err := conn.Exec(ctx, `
		UPDATE outbox
		SET processed = true, processed_at = $1
		WHERE id = $2
		`, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark message as processed: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return model.ErrMessageAlreadyProcessed
	}

	return nil
}
