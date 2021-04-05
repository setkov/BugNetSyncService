package BugNetService

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

// BugNet data service
type DataService struct {
	ConnectionString string
	Db               *sql.DB
}

// New data service
func NewDataService(connectionString string) *DataService{
	return &DataService{ConnectionString: connectionString}
}

// Open data connection
func (s *DataService) Open() error {
	var err error
	s.Db, err = sql.Open("sqlserver", s.ConnectionString)
	if err != nil {
		return err
	}
	return s.Db.Ping()
}

// Close data connection
func (s *DataService) Close() error {
	return s.Db.Close()
}

// Get message queue
func (s *DataService) GetMessageQueue() (*MessageQueue, error) {
	var que = MessageQueue{}

	rows, err := s.Db.Query("select top 10 * from dbo.Iserv_MessageQueue order by link desc")
	if err != nil {
		return &que, err
	}
	defer rows.Close()

	for rows.Next() {
		var mes Message
		err := rows.Scan(&mes.Link, &mes.Date, &mes.IssueId, &mes.TfsId, &mes.User, &mes.Operation, &mes.Message, &mes.DateSync)
		if err != nil {
			return &que, err
		}
		que.Messages = append(que.Messages, &mes)
	}
	return &que, nil
}

// Pull message for sync
func (s *DataService) PullMessage() (*Message, error) {
	var mes Message
	tsql := "select top 1 * from dbo.Iserv_MessageQueue where DateSync is null order by link"
	err := s.Db.QueryRow(tsql).Scan(&mes.Link, &mes.Date, &mes.IssueId, &mes.TfsId, &mes.User, &mes.Operation, &mes.Message, &mes.DateSync)
	return &mes, err
}

// Push message date sync
func (s *DataService) PushMessageDateSync(mes *Message) error {
	_, err := s.Db.Exec("update dbo.Iserv_MessageQueue set DateSync = @DateSync where link = @link", sql.Named("DateSync", mes.DateSync), sql.Named("link", mes.Link))
	if err != nil {
		return err
	}
	return nil
}
