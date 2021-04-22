package SyncService

import (
	"BugNetSyncService/BugNetService"
	"BugNetSyncService/TfsService"
	"fmt"
	"log"
	"time"
)

// Sync service
type SyncService struct {
	DataService *BugNetService.DataService
	TfsService  *TfsService.TfsService
	idleMode    bool
	stop        chan bool
}

// New sync service
func NewSyncService(b *BugNetService.DataService, t *TfsService.TfsService, idleMode bool) *SyncService {
	return &SyncService{
		DataService: b,
		TfsService:  t,
		stop:        make(chan bool),
		idleMode:    idleMode,
	}
}

// Start sync
func (s *SyncService) Start() {
	log.Print("Sync started.")
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			log.Print("Sync messages...")
			for {
				if err := s.syncMessage(); err != nil {
					log.Print(err)
					break
				}
				if s.idleMode {
					break
				}
			}
			log.Print("complited.")

			select {
			case <-ticker.C:
				continue
			case <-s.stop:
				log.Print("Sync stoped.")
				return
			}
		}
	}()
}

// Stop sync
func (s *SyncService) Stop() {
	s.stop <- true
}

// Sync message
func (s *SyncService) syncMessage() error {
	log.Print("PullMessage")
	message, err := s.DataService.PullMessage()
	if err != nil {
		return err
	} else {
		log.Print("Message: ", message)
	}
	comment := fmt.Sprintf("<p><i>%v %v %v:</i></p>%v", message.Date.Format("2006-01-02 15:04"), message.User.String, message.Operation.String, message.Message.String)
	log.Print("Comment: ", comment)

	log.Print("GetWorkItemsRelated")
	tfsIds := TfsService.TfsIds{Ids: []int{message.TfsId}}
	fields := []string{"System.WorkItemType", "System.State"}
	tfsWorkItems, err := s.TfsService.GetWorkItemsRelated(tfsIds, fields)
	if err != nil {
		return err
	} else {
		log.Print("Related WorkItems: ", tfsWorkItems)
	}

	for _, tfsWorkItem := range tfsWorkItems.Items {
		if (tfsWorkItem.Fields.WorkItemType == "Issue" || tfsWorkItem.Fields.WorkItemType == "Requirement" || tfsWorkItem.Fields.WorkItemType == "Bug") &&
			(tfsWorkItem.Fields.State == "Active" || tfsWorkItem.Fields.State == "Proposed" || tfsWorkItem.Fields.State == "Resolved") {
			log.Print("AddWorkItemComment")
			if s.idleMode {
				log.Print("IdleMode ON. Fake add comment to work item ", tfsWorkItem.Id)
			} else {
				workItem, err := s.TfsService.AddWorkItemComment(tfsWorkItem.Id, comment)
				if err != nil {
					return err
				} else {
					log.Print("Comment add to work item ", workItem.Id)
				}
			}
		}
	}

	log.Print("PushMessageDateSync")
	if s.idleMode {
		log.Print("IdleMode ON. Fake push message date sync ")
	} else {
		err = s.DataService.PushMessageDateSync(message)
		if err != nil {
			return err
		} else {
			log.Print("Accepted")
		}
	}
	return nil
}
