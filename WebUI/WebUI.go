package WebUI

import (
	"BugNetSyncService/BugNetService"
	"context"
	"log"
	"net/http"
)

//	WebUI
type WebUI struct {
	dataService *BugNetService.DataService
	webServer   *http.Server
}

// New WebUI
func NewWebUI(b *BugNetService.DataService) *WebUI {
	return &WebUI{
		dataService: b,
	}
}

// Start WebUI
func (w *WebUI) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", w.MessageQueueController)

	w.webServer = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Print("WebUI started.")
	go func() {
		w.webServer.ListenAndServe()
	}()
}

// Stop WebUI
func (w *WebUI) Stop() {
	log.Print("WebUI shutdown.")
	w.webServer.Shutdown(context.Background())
}
