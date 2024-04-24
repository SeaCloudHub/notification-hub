package postgrestore

import "time"

type BookQuerySchema struct {
	ISBN string `db:"isbn"`
	Name string `db:"name"`
}

type NotificationSchema struct {
	Id        int       `sql:"primary_key;auto_increment"`
	Uid       string    `sql:"type:VARCHAR(128);index"`
	From      string    `sql:"type:VARCHAR(128);not null"`
	To        string    `sql:"type:VARCHAR(128);not null"`
	Content   string    `sql:"type:VARCHAR(128)"`
	Status    string    `sql:"type:VARCHAR(128)"`
	CreatedAt time.Time `sql:"type:DATETIME;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `sql:"type:DATETIME;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (NotificationSchema) TableName() string {
	return "notifications"
}
