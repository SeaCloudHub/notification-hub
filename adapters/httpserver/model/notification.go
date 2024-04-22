package model

type Notification struct {
	UserId  string `json:"user_id"`
	Content string `json:"content"`
}

type NotificationRequest struct {
	From          string         `json:"from"`
	Notifications []Notification `json:"notifications"`
}
