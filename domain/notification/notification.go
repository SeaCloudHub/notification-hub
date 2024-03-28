package notification

import "time"

const (
	StatusReady      = "ready"
	StatusProcessing = "processing"
	StatusFailure    = "failure"
	StatusSuccess    = "success"
	StatusPending    = "pending"
)

type Storage interface {
}

type Book struct {
	ISBN      string
	Payload   string
	Status    string
	CreatedAt time.Time
	ReTriesNo int
}

func NewNotification(isbn string, payload string) Book {
	return Book{ISBN: isbn, Payload: payload, CreatedAt: time.Now(), ReTriesNo: 0, Status: StatusReady}
}
