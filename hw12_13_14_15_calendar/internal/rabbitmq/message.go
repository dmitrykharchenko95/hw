package rabbitmq

import "time"

type Notification struct {
	ID        int64
	UserID    int64
	Title     string
	StartDate time.Time
}

type Message struct {
	ContentType string
	Body        []byte
}
