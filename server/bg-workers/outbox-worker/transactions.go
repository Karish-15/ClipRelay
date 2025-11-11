package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type OutboxEvent struct {
	ID        int64
	EventType string
	Payload   string
}

type OutboxRepo struct {
	db           *sql.DB
	leaseTimeout time.Duration
}

func NewOutboxRepo(db *sql.DB, leaseTimeout time.Duration) *OutboxRepo {
	return &OutboxRepo{db: db, leaseTimeout: leaseTimeout}
}

func (r *OutboxRepo) FetchAndClaim(ctx context.Context, workerID string, limit int) ([]OutboxEvent, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, `
        SELECT id, event_type, payload
        FROM outbox
        WHERE processed = false
			AND (
                in_progress = false
                OR taken_at < (now() - $1::interval)
			)
        ORDER BY id
        LIMIT $2
        FOR UPDATE SKIP LOCKED
    `,
		fmt.Sprintf("%f seconds", r.leaseTimeout.Seconds()),
		limit,
	)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	var events []OutboxEvent
	for rows.Next() {
		var e OutboxEvent
		if err := rows.Scan(&e.ID, &e.EventType, &e.Payload); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		events = append(events, e)
	}

	if len(events) == 0 {
		_ = tx.Commit()
		return events, nil
	}

	ids := make([]int64, len(events))
	for i, e := range events {
		ids[i] = e.ID
	}

	_, err = tx.ExecContext(ctx, `
        UPDATE outbox
        SET in_progress = true,
            taken_at = now(),
            taken_by = $1
        WHERE id = ANY($2)
    `, workerID, pq.Array(ids))
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *OutboxRepo) MarkProcessed(ctx context.Context, ids []int64) error {
	_, err := r.db.ExecContext(ctx, `
        UPDATE outbox
        SET processed = true,
            in_progress = false,
            taken_by = NULL,
            taken_at = NULL
        WHERE id = ANY($1)
    `, pq.Array(ids))
	return err
}

func (r *OutboxRepo) Release(ctx context.Context, ids []int64) error {
	_, err := r.db.ExecContext(ctx, `
        UPDATE outbox
        SET in_progress = false,
            taken_by = NULL,
            taken_at = NULL
        WHERE id = ANY($1)
    `, pq.Array(ids))
	return err
}
