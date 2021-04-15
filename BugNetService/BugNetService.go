package BugNetService

import (
	"BugNetSyncService/Common"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

// BugNet data service
type DataService struct {
	ConnectionString string
	Db               *sql.DB
}

// New data service
func NewDataService(connectionString string) *DataService {
	return &DataService{ConnectionString: connectionString}
}

// Open data connection
func (s *DataService) Open() error {
	var err error
	s.Db, err = sql.Open("sqlserver", s.ConnectionString)
	if err != nil {
		return Common.NewError("Open sql connection. " + err.Error())
	}
	if err = s.Db.Ping(); err != nil {
		return Common.NewError("Test sql connection. " + err.Error())
	}
	return nil
}

// Close data connection
func (s *DataService) Close() error {
	if err := s.Db.Close(); err != nil {
		return Common.NewError("Close sql connection. " + err.Error())
	}
	return nil
}

// Get message queue
func (s *DataService) GetMessageQueue() (*MessageQueue, error) {
	var que = MessageQueue{}

	rows, err := s.Db.Query("select top 10 * from dbo.Iserv_MessageQueue order by link desc")
	if err != nil {
		return &que, Common.NewError("Get message queue. " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var mes Message
		if err := rows.Scan(&mes.Link, &mes.Date, &mes.IssueId, &mes.TfsId, &mes.User, &mes.Operation, &mes.Message, &mes.DateSync); err != nil {
			return &que, Common.NewError("Get message queue row. " + err.Error())
		}
		que.Messages = append(que.Messages, &mes)
	}
	return &que, nil
}

// Pull message for sync
func (s *DataService) PullMessage() (*Message, error) {
	var mes Message
	tsql := "select top 1 * from dbo.Iserv_MessageQueue where DateSync is null order by link"
	if err := s.Db.QueryRow(tsql).Scan(&mes.Link, &mes.Date, &mes.IssueId, &mes.TfsId, &mes.User, &mes.Operation, &mes.Message, &mes.DateSync); err != nil {
		if err == sql.ErrNoRows {
			return &mes, Common.NewWarning("Pull message. " + err.Error())
		} else {
			return &mes, Common.NewError("Pull message. " + err.Error())
		}
	}
	return &mes, nil
}

// Push message date sync
func (s *DataService) PushMessageDateSync(mes *Message) error {
	_, err := s.Db.Exec("update dbo.Iserv_MessageQueue set DateSync = GETDATE() where link = @link", sql.Named("link", mes.Link))
	if err != nil {
		return Common.NewError("Push message date sync. " + err.Error())
	}
	return nil
}
