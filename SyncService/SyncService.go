package SyncService

import (
	"BugNetSyncService/BugNetService"
	"BugNetSyncService/TfsService"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Sync service
type SyncService struct {
	DataService *BugNetService.DataService
	TfsService  *TfsService.TfsService
}

// New sync service
func NewSyncService(b *BugNetService.DataService, t *TfsService.TfsService) *SyncService {
	return &SyncService{DataService: b, TfsService: t}
}

// Sync message
func (s *SyncService) SyncMessage() error {
	// Pull message
	log.Print("PullMessage")
	message, err := s.DataService.PullMessage()
	if err != nil {
		return err
	} else {
		log.Print("Message: ", message.Link, message.Date, message.IssueId, message.TfsId, message.User, message.Operation)
	}
	comment := fmt.Sprintf("<p><i>%v %v %v:</i></p>%v", message.Date.Format("2006-01-02 15:04"), message.User.String, message.Operation.String, message.Message.String)
	log.Print("Comment: ", comment)

	// Get work items related
	log.Print("GetWorkItemsRelated")
	tfsIds := TfsService.TfsIds{Ids: []int{message.TfsId}}
	fields := []string{"System.WorkItemType", "System.State"}
	tfsWorkItems, err := s.TfsService.GetWorkItemsRelated(tfsIds, fields)
	if err != nil {
		return err
	} else {
		log.Print("Related WorkItems: ", tfsWorkItems)
	}

	// Sync work items
	for _, tfsWorkItem := range tfsWorkItems.Items {
		if (tfsWorkItem.Fields.WorkItemType == "Issue" || tfsWorkItem.Fields.WorkItemType == "Requirement" || tfsWorkItem.Fields.WorkItemType == "Bug") &&
			(tfsWorkItem.Fields.State == "Active" || tfsWorkItem.Fields.State == "Proposed" || tfsWorkItem.Fields.State == "Resolved") {
			// add work item comment
			log.Print("AddWorkItemComment")
			workItem, err := s.TfsService.AddWorkItemComment(tfsWorkItem.Id, comment)
			if err != nil {
				return err
			} else {
				log.Printf("Comment add to work item %v", workItem.Id)
			}
		}
	}

	//Push message date sync
	log.Print("PushMessageDateSync")
	_, offset := time.Now().Zone()
	loc := time.FixedZone("UTC", offset)
	message.DateSync = sql.NullTime{
		Time:  time.Now().In(loc),
		Valid: true,
	}
	err = s.DataService.PushMessageDateSync(message)
	if err != nil {
		return err
	} else {
		log.Print("DateSync: ", message.DateSync)
	}

	return nil
}
