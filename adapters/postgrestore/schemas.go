package postgrestore

import "time"

type BookQuerySchema struct {
	ISBN string `db:"isbn"`
	Name string `db:"name"`
}

type NotificationSchema struct {
	Uid       string    `gorm:"column:id"`
	From      string    `gorm:"column:from"`
	To        string    `gorm:"column:to"`
	Content   string    `gorm:"column:content"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (NotificationSchema) TableName() string {
	return "notifications"
}
