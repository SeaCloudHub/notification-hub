package postgrestore

import "time"

type BookQuerySchema struct {
	ISBN string `db:"isbn"`
	Name string `db:"name"`
}

type NotificationSchema struct {
	Uid       string    `gorm:"column:id"`
	From      string    `gorm:"column:from_user"`
	To        string    `gorm:"column:to_user"`
	Content   string    `gorm:"column:content"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	ViewedAt  time.Time `gorm:"column:viewed_at"`
}

func (NotificationSchema) TableName() string {
	return "notifications"
}
