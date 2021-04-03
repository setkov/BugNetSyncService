package main

import (
	"BugNetSyncService/BugNetService"
	"BugNetSyncService/ConfigService"
	"BugNetSyncService/TfsService"
	"fmt"
	"log"
)

func main() {
	log.Print("BugNetSyncService started.")

	log.Print("Load configuration...")
	var config ConfigService.Config
	err := config.Load()
	if err != nil {
		log.Fatal("Error: ", err.Error())
	} else {
		log.Print("Configuration loaded.")
	}

	log.Print("BugNet data service open...")
	bugNetService := BugNetService.DataService{ConnectionString: config.BugNetConnectionString}
	err = bugNetService.Open()
	if err != nil {
		log.Fatal("Error: ", err.Error())
	} else {
		log.Print("BugNet conected.")
	}

	// test getMessageQueue
	// log.Print("MessageQueue")
	// que, err := bugNetService.GetMessageQueue()
	// if err != nil {
	// 	log.Print("Error: ", err.Error())
	// } else {
	// 	for _, mes := range que.Messages {
	// 		log.Print(mes.Link, mes.Date, mes.IssueId, mes.TfsId, mes.User, mes.Operation, mes.DateSync)
	// 	}
	// }

	// test PullMessage
	log.Print("PullMessage")
	mes, err := bugNetService.PullMessage()
	if err != nil {
		log.Print("Error: ", err.Error())
	} else {
		log.Print(mes.Link, mes.Date, mes.IssueId, mes.TfsId, mes.User, mes.Operation, mes.DateSync)
	}

	// test GetWorkItemsRelated
	log.Print("GetWorkItemsRelated")
	tfsProvider := TfsService.NewTfsProvider(config.TfsBaseUri, config.TfsАuthorizationToken)
	tfsService := TfsService.NewTfsService(tfsProvider)
	//
	tfsIds := TfsService.TfsIds{Ids: []int{480565}}
	fields := []string{"System.WorkItemType", "System.State"}
	tfsWorkItems, err := tfsService.GetWorkItemsRelated(tfsIds, fields)
	if err != nil {
		log.Print("Error: ", err.Error())
	} else {
		log.Print("WorkItems: ", tfsWorkItems)
	}

	// // test AddWorkItemComment
	// log.Print("AddWorkItemComment")
	// comment := fmt.Sprintf("<p><i>%v %v %v:</i></p>%v", mes.Date.Format("2006-01-02 15:04"), mes.User.String, mes.Operation.String, mes.Message.String)
	// log.Print("Comment: ", comment)
	// tfsWorkItem, err := tfsService.AddWorkItemComment(483248, comment)
	// if err != nil {
	// 	log.Print("Error: ", err.Error())
	// } else {
	// 	log.Printf("Comment add to work item %v", tfsWorkItem.Id)
	// }

	// test PushMessageDateSync
	// log.Print("PushMessageDateSync")
	// _, offset := time.Now().Zone()
	// loc := time.FixedZone("UTC", offset)
	// mes.DateSync = sql.NullTime{
	// 	Time:  time.Now().In(loc),
	// 	Valid: true,
	// }
	// err = bugNetService.PushMessageDateSync(mes)
	// if err != nil {
	// 	log.Print("Error: ", err.Error())
	// } else {
	// 	log.Print(mes.Link, mes.Date, mes.IssueId, mes.TfsId, mes.User, mes.Operation, mes.DateSync)
	// }

	log.Print("BugNetSyncService stoped.")
}
