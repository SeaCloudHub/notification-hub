package postgrestore

import (
	"context"

	"github.com/SeaCloudHub/notification-hub/domain/notification"
	"gorm.io/gorm"
)

type NotificationStore struct {
	db *gorm.DB
}

func NewNotificationStore(db *gorm.DB) *NotificationStore {
	return &NotificationStore{db}
}

func (s *NotificationStore) Create(ctx context.Context, notification *notification.Notification) error {
	notiSchema := NotificationSchema{
		Uid:     notification.Uid,
		From:    notification.From,
		To:      notification.To,
		Content: notification.Content,
		Status:  notification.Status,
	}

	return s.db.WithContext(ctx).Create(&notiSchema).Error
}
