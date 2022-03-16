package main

import (
	"BugNetSyncService/BugNetService"
	"BugNetSyncService/Common"
	"BugNetSyncService/SyncService"
	"BugNetSyncService/TfsService"
	"BugNetSyncService/WebUI"
	"fmt"
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

	msTeamsServise := Common.NewMSTeamsService(config.MSTeamsWebhookUrl, "BugNetSyncService")
	if !config.IdleMode {
		msTeamsServise.SendMessage("Service started.")
	}

	log.Print("BugNet data service open...")
	bugNetService, err := BugNetService.NewDataService(config)
	if err != nil {
		log.Print(err)
	} else {
		log.Print("conected.")
	}

	webUI := WebUI.NewWebUI(bugNetService)
	webUI.Start()

	tfsProvider := TfsService.NewTfsProvider(config.TfsBaseUri, config.TfsАuthorizationToken)
	tfsService := TfsService.NewTfsService(tfsProvider)

	syncService := SyncService.NewSyncService(bugNetService, tfsService, msTeamsServise, config.IdleMode)
	syncService.Start()

	/* 	// test Attach file
	   	time.Sleep(2 * time.Second)
	   	log.Print("test Attach file")
	   	mes, err := bugNetService.GetMessage(23115)
	   	if err != nil {
	   		log.Print(err)
	   	}
	   	log.Print(mes)
	   	bytes, err := bugNetService.LoadAttachment(int(mes.AttachmentId.Int32))
	   	if err != nil {
	   		log.Print(err)
	   	}
	   	//log.Print(bytes)
	   	workItem, err := tfsService.AddWorkItemAttachment(mes.TfsId, mes.FileName.String, bytes)
	   	//workItem, err := tfsService.AddWorkItemAttachment(290704, mes.FileName.String, bytes)
	   	if err != nil {
	   		log.Print(err)
	   	} else {
	   		log.Print("Loaded to work item ", workItem.Id)
	   	} */
	// test message image
	message, err := bugNetService.GetMessage(41104)
	if err == nil {
		messageImages := BugNetService.GetMessageImages(message.Message.String)
		for _, image := range messageImages.Images {
			imageSrc := BugNetService.GetImageSrc(image)
			fileName := fmt.Sprintf("d:\\%s.%s", imageSrc.Name, imageSrc.Ext)
			log.Print(fileName)
			imageSrc.SaveAsFile(fileName)
		}
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	webUI.Stop()
	syncService.Stop()
	// wait services stoped
	time.Sleep(1 * time.Second)

	if !config.IdleMode {
		msTeamsServise.SendMessage("Service stoped.")
	}
	log.Print("BugNetSyncService stoped.")
}
