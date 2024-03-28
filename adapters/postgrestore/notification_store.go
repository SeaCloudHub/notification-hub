package postgrestore

import (
	"context"
	"fmt"

	"github.com/SeaCloudHub/notification-hub/domain/book"
	"github.com/jmoiron/sqlx"
)

type NotificationStore struct {
	db *sqlx.DB
}

func NewNotificationStore(db *sqlx.DB) *NotificationStore {
	return &NotificationStore{db}
}

func (s *NotificationStore) Save(ctx context.Context, b *book.Book) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO books(isbn,name) VALUES ($1,$2)`, b.ISBN, b.Name)
	if err != nil {
		return fmt.Errorf("cannot save the book: %w", err)
	}
	return nil
}
