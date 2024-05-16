package model

import (
	"context"

	"github.com/SeaCloudHub/notification-hub/domain/notification"
	"github.com/SeaCloudHub/notification-hub/pkg/pagination"
	"github.com/SeaCloudHub/notification-hub/pkg/validation"
)

type UpdateViewedTimeRequest struct {
	IdNotification string `json:"id_noti"`
}

type MarkEntireViewedRequest struct {
}

type MarkEntireViewedResponse struct {
	Status string `json:"status"`
}

type UpdateViewedTimeResponse struct {
	Status string `json:"status"`
}

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
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Cursor string `query:"cursor" validate:"omitempty,base64url"`
}

type ListPageEntriesRequest struct {
	Page  int `query:"page" validate:"required,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=100"`
}

func (r *ListPageEntriesRequest) Validate(ctx context.Context) error {
	if r.Limit == 0 {
		r.Limit = 10
	}

	if r.Page == 0 {
		r.Page = 1
	}

	return validation.Validate().StructCtx(ctx, r)
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

type ListPageEntriesResponse struct {
	Entries    []*notification.Notification `json:"entries"`
	Pagination pagination.PageInfo          `json:"pagination"`
}
