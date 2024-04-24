package postgrestore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SeaCloudHub/notification-hub/domain/identity"
	"github.com/SeaCloudHub/notification-hub/domain/notification"
	"github.com/SeaCloudHub/notification-hub/pkg/pagination"

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

func (s *NotificationStore) UpdateStatusByUid(ctx context.Context, uid string, status string) error {
	return s.db.WithContext(ctx).Model(&NotificationSchema{}).
		Where("id = ?", uid).Update("status", status).Error
}

func (s *NotificationStore) GetByUid(ctx context.Context, uid string) (*notification.Notification, error) {
	var noti NotificationSchema
	err := s.db.WithContext(ctx).Where("id = ?", uid).First(&noti).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, identity.ErrIdentityNotFound
		}

		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	return &notification.Notification{
		Uid:     noti.Uid,
		From:    noti.From,
		To:      noti.To,
		Content: noti.Content,
		Status:  noti.Status,
	}, nil
}

func (s *NotificationStore) ListByUserId(ctx context.Context, userId string) ([]*notification.Notification, error) {
	var notis []NotificationSchema
	if err := s.db.WithContext(ctx).Where("to = ?", userId).Find(&notis).Error; err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	notisResult := make([]*notification.Notification, 0, len(notis))
	for _, notiSchema := range notis {
		notisResult = append(notisResult, &notification.Notification{
			Uid:     notiSchema.Uid,
			From:    notiSchema.From,
			To:      notiSchema.To,
			Content: notiSchema.Content,
			Status:  notiSchema.Status,
		})
	}

	return notisResult, nil
}

func (s *NotificationStore) ListByUserIdUsingCursor(ctx context.Context, userId string, cursor *pagination.Cursor) ([]*notification.Notification, error) {
	var notiSchemas []NotificationSchema

	// parse cursor
	cursorObj, err := pagination.DecodeToken[fsCursor](cursor.Token)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", identity.ErrIdentityInvalidCursor, err)
	}

	query := s.db.WithContext(ctx).Where("to = ?", userId)
	if cursorObj.CreatedAt != nil {
		query = query.Where("created_at >= ?", cursorObj.CreatedAt)
	}

	if err := query.Limit(cursor.Limit + 1).Find(&notiSchemas).Error; err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	if len(notiSchemas) > cursor.Limit {
		cursor.SetNextToken(pagination.EncodeToken(fsCursor{CreatedAt: &notiSchemas[cursor.Limit].CreatedAt}))
		notiSchemas = notiSchemas[:cursor.Limit]
	}

	notisResult := make([]*notification.Notification, 0, len(notiSchemas))
	for _, notiSchema := range notiSchemas {
		notisResult = append(notisResult, &notification.Notification{
			Uid:     notiSchema.Uid,
			From:    notiSchema.From,
			To:      notiSchema.To,
			Content: notiSchema.Content,
			Status:  notiSchema.Status,
		})
	}

	return notisResult, nil

}

type fsCursor struct {
	CreatedAt *time.Time
}
