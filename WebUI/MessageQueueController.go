package WebUI

import (
	"fmt"
	"net/http"
)

func (webUI *WebUI) MessageQueueController(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "MessageQueue")
}
