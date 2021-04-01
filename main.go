package main

import (
	"BugNetSyncService/BugNetService"
	"BugNetSyncService/ConfigService"
	"BugNetSyncService/TfsService"
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
	tfsService := TfsService.NewTfsService(config.TfsBaseUri, config.TfsАuthorizationToken)
	tfsProvider := TfsService.NewTfsProvider(&tfsService)
	//
	tfsIds := TfsService.TfsIds{Ids: []int{480565}}
	fields := []string{"System.WorkItemType", "System.State"}
	tfsWorkItems, err := tfsProvider.GetWorkItemsRelated(tfsIds, fields)
	if err != nil {
		log.Print("Error: ", err.Error())
	} else {
		log.Print("WorkItems: ", tfsWorkItems)
	}

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
