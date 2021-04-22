package BugNetService

import (
	"database/sql"
	"time"
)

type Message struct {
	Id           int
	Link         int
	Date         time.Time
	IssueId      int
	TfsId        int
	User         sql.NullString
	Operation    sql.NullString
	Message      sql.NullString
	DateSync     sql.NullTime
	IssueUrl     sql.NullString
	TfsUrl       sql.NullString
	AttachmentId sql.NullInt32
	FileName     sql.NullString
	ContentType  sql.NullString
	FileUrl      sql.NullString
}

type MessageQueue struct {
	Messages []*Message
}
