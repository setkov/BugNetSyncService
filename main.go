package main

import (
	"BugNetSyncService/BugNetService"
	"BugNetSyncService/ConfigService"
	"BugNetSyncService/SyncService"
	"BugNetSyncService/TfsService"
	"fmt"
	"log"
	"time"
)

func main() {
	log.Print("BugNetSyncService started.")

	log.Print("Load configuration...")
	config := ConfigService.NewConfigService()
	log.Print("loaded.")

	log.Print("BugNet data service open...")
	bugNetService := BugNetService.NewDataService(config.BugNetConnectionString)
	if err := bugNetService.Open(); err != nil {
		log.Fatal("Error: ", err.Error())
	} else {
		log.Print("conected.")
	}

	tfsProvider := TfsService.NewTfsProvider(config.TfsBaseUri, config.TfsАuthorizationToken)
	tfsService := TfsService.NewTfsService(tfsProvider)

	syncService := SyncService.NewSyncService(bugNetService, tfsService)
	syncService.Start()

	fmt.Scanln()
	syncService.Stop()
	// wait sync service stoped
	time.Sleep(1 * time.Second)

	log.Print("BugNetSyncService stoped.")
}
