package WebUI

import (
	"BugNetSyncService/BugNetService"
	"context"
	"html/template"
	"log"
	"net/http"
)

//	WebUI
type WebUI struct {
	dataService *BugNetService.DataService
	server      *http.Server
	funcMap     template.FuncMap
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

	w.server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	w.funcMap = template.FuncMap{
		"html":     TextToHtml,
		"text":     HtmlToText,
		"trim":     TrimText,
		"datetime": DateTime,
	}

	log.Print("WebUI started.")
	go func() {
		w.server.ListenAndServe()
	}()
}

// Stop WebUI
func (w *WebUI) Stop() {
	log.Print("WebUI shutdown.")
	w.server.Shutdown(context.Background())
}
