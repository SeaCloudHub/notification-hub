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

type Storage interface {
	Create(ctx context.Context, notification *Notification) (*Notification, error)
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
	return Notification{Uid: id, From: from, Status: StatusReady, Content: content}
}
