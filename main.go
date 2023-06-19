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
	log.Print("Load configuration.")
	config, err := Common.NewConfig()
	if err != nil {
		log.Print(err)
	}

	//messengerServise := Common.NewMSTeamsService(config.MSTeamsWebhookUrl, "BugNetSyncService")
	messengerServise := Common.NewTelegramService(config.TelegramToken, config.TelegramChatId, "BugNetSyncService")

	var message = "Service started"
	if config.IdleMode {
		message += " in idle mode"
	}
	message += " - ver: " + config.ApplicationVersion
	log.Print(message)
	if !config.IdleMode {
		if err := messengerServise.SendMessage(message); err != nil {
			log.Print(err)
		}
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

	tfsProvider := TfsService.NewTfsProvider(config.TfsBaseUri, config.TfsAuthorizationToken)
	tfsService := TfsService.NewTfsService(tfsProvider)

	syncService := SyncService.NewSyncService(bugNetService, tfsService, messengerServise, config.IdleMode)
	syncService.Start()

	/* 	// test sync message
	syncService := SyncService.NewSyncService(bugNetService, tfsService, msTeamsServise, false) // Set Idle Mode
	message, err := bugNetService.GetMessage(41196)
	if err != nil {
		log.Print(err)
	} else {
		err = syncService.SyncMessage(message)
		if err != nil {
			log.Print(err)
		}
	} */
	/* 	// test message image
	message, err := bugNetService.GetMessage(41196)
	if err != nil {
		log.Print(err)
	} else {
		//log.Print(message.Message.String)
		messageImages := BugNetService.GetMessageImages(message.Message.String)
		for _, image := range messageImages.Images {
			fileName := fmt.Sprintf("d:\\%s.%s", image.ImageSrc.Name, image.ImageSrc.Ext)
			log.Print(fileName)
			err := image.ImageSrc.SaveAsFile(fileName)
			if err != nil {
				log.Print(err)
			}
		}
	} */
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

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	webUI.Stop()
	syncService.Stop()
	// wait services stoped
	time.Sleep(1 * time.Second)

	message = "Service stoped"
	if config.IdleMode {
		message += " in idle mode"
	}
	if !config.IdleMode {
		if err := messengerServise.SendMessage(message); err != nil {
			log.Print(err)
		}
	}
	log.Print(message)
}
