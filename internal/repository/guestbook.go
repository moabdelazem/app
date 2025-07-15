package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/moabdelazem/app/internal/database"
	"github.com/moabdelazem/app/internal/models"
)

type GuestBookRepository struct {
	db *database.DB
}

func NewGuestBookRepository(db *database.DB) *GuestBookRepository {
	return &GuestBookRepository{db: db}
}

func (r *GuestBookRepository) CreateTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS guest_book_messages (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
		
		CREATE INDEX IF NOT EXISTS idx_guest_book_created_at ON guest_book_messages(created_at DESC);
	`

	_, err := r.db.Pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create guest_book_messages table: %w", err)
	}

	return nil
}

func (r *GuestBookRepository) Create(ctx context.Context, msg *models.CreateGuestBookMessage) (*models.GuestBookMessage, error) {
	query := `
		INSERT INTO guest_book_messages (name, email, message)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, message, created_at, updated_at
	`

	var result models.GuestBookMessage
	err := r.db.Pool.QueryRow(ctx, query, msg.Name, msg.Email, msg.Message).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.Message,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create guest book message: %w", err)
	}

	return &result, nil
}

func (r *GuestBookRepository) GetAll(ctx context.Context, limit, offset int) ([]models.GuestBookMessage, error) {
	query := `
		SELECT id, name, email, message, created_at, updated_at
		FROM guest_book_messages
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get guest book messages: %w", err)
	}
	defer rows.Close()

	var messages []models.GuestBookMessage
	for rows.Next() {
		var msg models.GuestBookMessage
		err := rows.Scan(
			&msg.ID,
			&msg.Name,
			&msg.Email,
			&msg.Message,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan guest book message: %w", err)
		}
		messages = append(messages, msg)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating guest book messages: %w", rows.Err())
	}

	return messages, nil
}

func (r *GuestBookRepository) GetByID(ctx context.Context, id int) (*models.GuestBookMessage, error) {
	query := `
		SELECT id, name, email, message, created_at, updated_at
		FROM guest_book_messages
		WHERE id = $1
	`

	var msg models.GuestBookMessage
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&msg.ID,
		&msg.Name,
		&msg.Email,
		&msg.Message,
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("guest book message not found")
		}
		return nil, fmt.Errorf("failed to get guest book message: %w", err)
	}

	return &msg, nil
}

func (r *GuestBookRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM guest_book_messages`

	var count int
	err := r.db.Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count guest book messages: %w", err)
	}

	return count, nil
}
