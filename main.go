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
	bugNetService := BugNetService.NewDataService(config.BugNetConnectionString, config.BugNetDomainUrl, config.BugNetАuthorizationToken)
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

	// test send message
	msTeamsServise := Common.NewMSTeamsService("https://outlook.office.com/webhook/3e99495e-14ea-4806-ba60-616ee03e4d1e@9282aa4b-7330-45c1-8391-ef6e504d84b9/IncomingWebhook/b356823fa9bd452ea09c547d1af5040c/6bc06d46-a37e-49d2-b9b4-3157037ada58", "BugNetSyncService")
	msTeamsServise.SendMessage("test")

	/* 	// test Attach file
	   	time.Sleep(2 * time.Second)
	   	log.Print("test Attach file")
	   	mes, err := bugNetService.GetMessage(1301)
	   	if err != nil {
	   		log.Print(err)
	   	}
	   	log.Print(mes)
	   	bytes, err := bugNetService.LoadAttachment(mes)
	   	if err != nil {
	   		log.Print(err)
	   	}
	   	//log.Print(bytes)
	   	// workItem, err := tfsService.AddWorkItemAttachment(mes.TfsId, mes.FileName.String, bytes)
	   	workItem, err := tfsService.AddWorkItemAttachment(290704, mes.FileName.String, bytes)
	   	if err != nil {
	   		log.Print(err)
	   	} else {
	   		log.Print("Loaded to work item ", workItem.Id)
	   	} */

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	webUI.Stop()
	syncService.Stop()
	// wait services stoped
	time.Sleep(1 * time.Second)

	log.Print("BugNetSyncService stoped.")
}
