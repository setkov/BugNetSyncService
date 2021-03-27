package BugNetService

import (
	"database/sql"
	"time"
)

type Message struct {
	Link      int
	Date      time.Time
	IssueId   int
	TfsId     int
	User      sql.NullString
	Operation sql.NullString
	Message   sql.NullString
	DateSync  sql.NullTime
}

type MessageQueue struct {
	Messages []*Message
}
