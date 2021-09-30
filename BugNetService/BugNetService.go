package BugNetService

import (
	"BugNetSyncService/Common"
	"database/sql"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
)

// BugNet data service
type DataService struct {
	AttachmentServiceUrl string
	db                   *sql.DB
}

// New data service
func NewDataService(config *Common.Config) (*DataService, error) {
	var dataService = &DataService{
		AttachmentServiceUrl: config.AttachmentServiceUrl,
	}

	// open data connection
	var err error
	dataService.db, err = sql.Open("sqlserver", config.ConnectionString)
	if err != nil {
		return dataService, Common.NewError("Open sql connection. " + err.Error())
	}
	if err = dataService.db.Ping(); err != nil {
		return dataService, Common.NewError("Test sql connection. " + err.Error())
	}

	return dataService, nil
}

// Close data connection
func (s *DataService) Close() error {
	if err := s.db.Close(); err != nil {
		return Common.NewError("Close sql connection. " + err.Error())
	}
	return nil
}

// Get message queue
func (s *DataService) GetMessageQueue(top int) (*MessageQueue, error) {
	var que = MessageQueue{}

	rows, err := s.db.Query("select top (@top) [Id],[Date],[IssueId],[TfsId],[User],[Operation],[Message],[DateSync],[IssueUrl],[TfsUrl],[AttachmentId],[FileName],[ContentType],[FileUrl] from dbo.Iserv_MessageQueue order by id desc", sql.Named("top", top))
	if err != nil {
		return &que, Common.NewError("Get message queue. " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var mes Message
		if err := rows.Scan(&mes.Id, &mes.Date, &mes.IssueId, &mes.TfsId, &mes.User, &mes.Operation, &mes.Message, &mes.DateSync, &mes.IssueUrl, &mes.TfsUrl, &mes.AttachmentId, &mes.FileName, &mes.ContentType, &mes.FileUrl); err != nil {
			return &que, Common.NewError("Get message queue row. " + err.Error())
		}
		que.Messages = append(que.Messages, &mes)
	}
	return &que, nil
}

// Get message by id
func (s *DataService) GetMessage(id int) (*Message, error) {
	var mes Message
	tsql := "select [Id],[Date],[IssueId],[TfsId],[User],[Operation],[Message],[DateSync],[IssueUrl],[TfsUrl],[AttachmentId],[FileName],[ContentType],[FileUrl] from dbo.Iserv_MessageQueue where Id = @Id"
	if err := s.db.QueryRow(tsql, sql.Named("Id", id)).Scan(&mes.Id, &mes.Date, &mes.IssueId, &mes.TfsId, &mes.User, &mes.Operation, &mes.Message, &mes.DateSync, &mes.IssueUrl, &mes.TfsUrl, &mes.AttachmentId, &mes.FileName, &mes.ContentType, &mes.FileUrl); err != nil {
		if err == sql.ErrNoRows {
			return &mes, Common.NewWarning("Pull message. " + err.Error())
		} else {
			return &mes, Common.NewError("Pull message. " + err.Error())
		}
	}
	return &mes, nil
}

// Pull message for sync
func (s *DataService) PullMessage() (*Message, error) {
	var mes Message
	tsql := "select top 1 [Id],[Date],[IssueId],[TfsId],[User],[Operation],[Message],[DateSync],[IssueUrl],[TfsUrl],[AttachmentId],[FileName],[ContentType],[FileUrl] from dbo.Iserv_MessageQueue where DateSync is null order by id"
	if err := s.db.QueryRow(tsql).Scan(&mes.Id, &mes.Date, &mes.IssueId, &mes.TfsId, &mes.User, &mes.Operation, &mes.Message, &mes.DateSync, &mes.IssueUrl, &mes.TfsUrl, &mes.AttachmentId, &mes.FileName, &mes.ContentType, &mes.FileUrl); err != nil {
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
	_, err := s.db.Exec("update dbo.Iserv_MessageQueue set DateSync = GETDATE() where Id = @Id", sql.Named("Id", mes.Id))
	if err != nil {
		return Common.NewError("Push message date sync. " + err.Error())
	}
	return nil
}

// Load attachment
func (s *DataService) LoadAttachment(id int) ([]byte, error) {
	var bytes []byte

	client := &http.Client{}
	req, err := http.NewRequest("GET", s.AttachmentServiceUrl+"?id="+strconv.Itoa(id), nil)
	if err != nil {
		return bytes, Common.NewError("New request. " + err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		return bytes, Common.NewError("Do request. " + err.Error())
	}
	defer resp.Body.Close()

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return bytes, Common.NewError("Read request. " + err.Error())
	}

	return bytes, nil
}
