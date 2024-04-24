package notification

import (
	"context"
)

const (
	StatusReady      = "ready"
	StatusProcessing = "processing"
	StatusFailure    = "failure"
	StatusSuccess    = "success"
	StatusPending    = "pending"
)

type Store interface {
	Create(ctx context.Context, notification *Notification) error
	UpdateStatusByUid(ctx context.Context, uid string, status string) error
	GetByUid(ctx context.Context, uid string) (*Notification, error)
	ListByUserId(ctx context.Context, userId string) ([]*Notification, error)
}

type Notification struct {
	Uid     string
	From    string
	To      string
	Content string
	Status  string
}

type SetOfNotifications struct {
	Noitications []*Notification
}

func NewNotification(id string, from string, to string, content string) Notification {
	return Notification{Uid: id, From: from, Status: StatusReady, Content: content, To: to}
}
