package main

import (
	"BugNetSyncService/BugNetService"
	"BugNetSyncService/ConfigService"
	"BugNetSyncService/TfsService"
	"log"
	"net/http"
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

	// test GetWorkItemsRelations
	log.Print("GetWorkItemsRelations")
	tfsService := TfsService.TfsService{BaseUri: config.TfsBaseUri, АuthorizationToken: config.TfsАuthorizationToken, Client: &http.Client{}}
	tfsIds := TfsService.TfsIds{Ids: []int{480565}}
	var rel TfsService.TfsRelations
	err = rel.Load(&tfsService, tfsIds, "System.LinkTypes.Related")
	if err != nil {
		log.Print("Error: ", err.Error())
	} else {
		// log.Print(res)
		log.Print("Related: ", rel)

		tfsIds.AddTargets(rel)
		log.Print("tfsIds: ", tfsIds)

		// get child relations
		var rel TfsService.TfsRelations
		err = rel.Load(&tfsService, tfsIds, "System.LinkTypes.Hierarchy-Forward")
		if err != nil {
			log.Print("Error: ", err.Error())
		} else {
			log.Print("Child: ", rel)

			tfsIds.AddTargets(rel)
			log.Print("tfsIds: ", tfsIds)

			// get tfs work items
			fields := []string{"System.WorkItemType", "System.State"}
			var items TfsService.TfsWorkItems
			err = items.Load(&tfsService, tfsIds, fields)
			if err != nil {
				log.Print("Error: ", err.Error())
			} else {
				log.Print("WorkItems: ", items)
			}
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
