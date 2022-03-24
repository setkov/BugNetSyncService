package SyncService

import (
	"BugNetSyncService/BugNetService"
	"BugNetSyncService/Common"
	"BugNetSyncService/TfsService"
	"fmt"
	"log"
	"time"
)

// Sync service
type SyncService struct {
	DataService    *BugNetService.DataService
	TfsService     *TfsService.TfsService
	MSTeamsService *Common.MSTeamsService
	idleMode       bool
	stop           chan bool
}

// New sync service
func NewSyncService(b *BugNetService.DataService, t *TfsService.TfsService, s *Common.MSTeamsService, idleMode bool) *SyncService {
	return &SyncService{
		DataService:    b,
		TfsService:     t,
		MSTeamsService: s,
		idleMode:       idleMode,
		stop:           make(chan bool),
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
				log.Print("PullMessage")
				message, err := s.DataService.PullMessage()
				if err != nil {
					LogError(s, err)
					break
				} else {
					log.Print("Message: ", message)
				}
				err = s.SyncMessage(message)
				if err != nil {
					LogError(s, err)
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

// Log error
func LogError(s *SyncService, err error) {
	log.Print(err)
	errorWithCategory, ok := err.(*Common.ErrorWithCategory)
	if ok {
		if !s.idleMode && errorWithCategory.Category() == Common.Error {
			s.MSTeamsService.SendMessage(errorWithCategory.Error())
		}
	}
}

// Sync message
func (s *SyncService) SyncMessage(message *BugNetService.Message) error {
	comment := fmt.Sprintf("<p><i>%v %v %v:</i></p>%v", message.Date.Format("2006-01-02 15:04"), message.User.String, message.Operation.String, message.Message.String)

	// find images in comment
	images := BugNetService.GetMessageImages(comment)
	// reverse loop images
	for i := len(images.Images) - 1; i >= 0; i-- {
		image := images.Images[i]

		// decode image
		imageBody, err := image.ImageSrc.DecodeBody()
		if err != nil {
			return Common.NewError("Decode image from comment. " + err.Error())
		}

		imageName := fmt.Sprintf("%v.%v", image.ImageSrc.Name, image.ImageSrc.Ext)
		log.Printf("Image %v: %v", i, imageName)

		// add images to TFS
		var ref TfsService.AttachmentReference
		if s.idleMode {
			log.Print("IdleMode ON. Fake added commeent image to work item ", message.TfsId)
		} else {
			ref, err = s.TfsService.Provider.CreateAttachment(imageName, imageBody)
			if err != nil {
				return Common.NewError("Create image attachment. " + err.Error())
			}
		}
		// replace images in comment
		comment = fmt.Sprintf("%v <img src=\"%v\" alt=\"%v.%v\" /> %v", comment[:image.StartPosition], ref.Url, image.ImageSrc.Name, image.ImageSrc.Ext, comment[image.StopPosition:])
	}
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

	// sync comment
	for _, tfsWorkItem := range tfsWorkItems.Items {
		if tfsWorkItem.Id == message.TfsId ||
			(tfsWorkItem.Fields.WorkItemType == "Issue" || tfsWorkItem.Fields.WorkItemType == "Requirement" || tfsWorkItem.Fields.WorkItemType == "Bug") &&
				(tfsWorkItem.Fields.State == "Active" || tfsWorkItem.Fields.State == "Proposed" || tfsWorkItem.Fields.State == "Resolved") {
			log.Print("AddWorkItemComment")
			if s.idleMode {
				log.Print("IdleMode ON. Fake added to work item ", tfsWorkItem.Id)
			} else {
				workItem, err := s.TfsService.AddWorkItemComment(tfsWorkItem.Id, comment)
				if err != nil {
					return err
				} else {
					log.Print("Added to work item ", workItem.Id)
				}
			}
		}
	}

	// sync attachment
	if message.Operation.String == "add attachment" {
		log.Print("AddWorkItemAttachment")
		if s.idleMode {
			log.Print("IdleMode ON. Fake added to work item ", message.TfsId)
		} else {
			log.Print("LoadAttachment")
			bytes, err := s.DataService.LoadAttachment(int(message.AttachmentId.Int32))
			if err != nil {
				return err
			}
			workItem, err := s.TfsService.AddWorkItemAttachment(message.TfsId, message.FileName.String, bytes)
			if err != nil {
				return err
			} else {
				log.Print("Loaded to work item ", workItem.Id)
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
