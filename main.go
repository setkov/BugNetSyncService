package main

import (
	"BugNetSyncService/BugNetService"
	"BugNetSyncService/Common"
	"BugNetSyncService/SyncService"
	"BugNetSyncService/TfsService"
	"BugNetSyncService/WebUI"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Print("BugNetSyncService started.")

	log.Print("Load configuration.")
	config, err := Common.NewConfig()
	if err != nil {
		log.Print(err)
	}

	log.Print("BugNet data service open...")
	bugNetService := BugNetService.NewDataService(config.BugNetConnectionString)
	if err := bugNetService.Open(); err != nil {
		log.Print(err)
	} else {
		log.Print("conected.")
	}

	webUI := WebUI.NewWebUI(bugNetService)
	webUI.Start()

	tfsProvider := TfsService.NewTfsProvider(config.TfsBaseUri, config.TfsАuthorizationToken)
	tfsService := TfsService.NewTfsService(tfsProvider)

	syncService := SyncService.NewSyncService(bugNetService, tfsService, config.IdleMode)
	syncService.Start()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	webUI.Stop()
	syncService.Stop()
	// wait services stoped
	time.Sleep(1 * time.Second)

	log.Print("BugNetSyncService stoped.")
}
