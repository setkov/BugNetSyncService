package main

import (
	"BugNetSyncService/BugNetService"
	"BugNetSyncService/ConfigService"
	"BugNetSyncService/SyncService"
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

	tfsProvider := TfsService.NewTfsProvider(config.TfsBaseUri, config.TfsАuthorizationToken)
	tfsService := TfsService.NewTfsService(tfsProvider)
	syncService := SyncService.NewSyncService(&bugNetService, tfsService)

	log.Print("Sync message...")
	err = syncService.SyncMessage()
	if err != nil {
		log.Print("Error: ", err.Error())
	} else {
		log.Print("Complited.")
	}

	log.Print("BugNetSyncService stoped.")
}
