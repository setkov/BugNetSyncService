package main

import (
	"BugNetSyncService/BugNetService"
	"BugNetSyncService/TfsService"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	var err error
	log.Print("BugNetSyncService started.")

	log.Print("Load configuration...")
	var config Config = Config{}
	err = config.Load()
	if err != nil {
		log.Fatal("Error: ", err.Error())
	} else {
		log.Print("Configuration loaded.")
	}

	log.Print("BugNet data service open...")
	var bugNetService BugNetService.DataService = BugNetService.DataService{ConnectionString: config.BugNetConnectionString}
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

	// test GetWorkItemsRelations
	log.Print("GetWorkItemsRelations")
	var tfsService TfsService.TfsService = TfsService.TfsService{BaseUri: config.TfsBaseUri, АuthorizationToken: config.TfsАuthorizationToken, Client: &http.Client{}}
	res, err := tfsService.GetWorkItemsRelations([]int{480565}, "System.LinkTypes.Related")
	if err != nil {
		log.Print("Error: ", err.Error())
	} else {
		// log.Print(res)
		relations := TfsService.TfsRelations{}
		err := json.Unmarshal([]byte(res), &relations)
		if err != nil {
			log.Print("Error: ", err.Error())
		} else {
			log.Print(relations)
		}
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
