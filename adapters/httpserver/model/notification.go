package model

import (
	"context"

	"github.com/SeaCloudHub/notification-hub/domain/notification"
	"github.com/SeaCloudHub/notification-hub/pkg/validation"
)

type Notification struct {
	UserId  string `json:"user_id"`
	Content string `json:"content"`
}

type NotificationRequest struct {
	From          string         `json:"from"`
	Notifications []Notification `json:"notifications"`
}

type NotificationResponse struct {
	Status string `json:"status"`
}

type ListEntriesRequest struct {
	ID     string `param:"id" validate:"required,uuid" swaggerignore:"true"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Cursor string `query:"cursor" validate:"omitempty,base64url"`
}

func (r *ListEntriesRequest) Validate(ctx context.Context) error {
	if r.Limit <= 0 {
		r.Limit = 10
	}

	return validation.Validate().StructCtx(ctx, r)
}

type ListEntriesResponse struct {
	Entries []*notification.Notification `json:"entries"`
	Cursor  string                       `json:"cursor"`
}
