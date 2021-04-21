package WebUI

import (
	"BugNetSyncService/Common"
	"fmt"
	"html/template"
	"net/http"
)

func (webUI *WebUI) MessageQueueController(w http.ResponseWriter, r *http.Request) {
	messageQueue, err := webUI.dataService.GetMessageQueue(100)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		tmpl, err := template.New("MessageQueueView.html").Funcs(webUI.funcMap).ParseFiles("WebUI/MessageQueueView.html")
		if err != nil {
			fmt.Fprint(w, Common.NewError("Parse temptate MessageQueueView. "+err.Error()))
		} else {
			if err := tmpl.Execute(w, messageQueue); err != nil {
				fmt.Fprint(w, Common.NewError("Execute temptate MessageQueueView. "+err.Error()))
			}
		}
	}
}
