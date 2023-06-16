package Common

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type MSTeamsService struct {
	WebhookUrl string
	Sender     string
}

// New MSTeamsService
func NewMSTeamsService(webhookUrl string, sender string) MSTeamsService {
	return MSTeamsService{
		WebhookUrl: webhookUrl,
		Sender:     sender,
	}
}

// Send message
func (s MSTeamsService) SendMessage(message string) error {
	body, err := json.Marshal(map[string]interface{}{"title": s.Sender, "text": message})
	if err != nil {
		return NewError("Prepare MSTeams message. " + err.Error())
	}

	resp, err := http.Post(s.WebhookUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return NewError("Send MSTeams message. " + err.Error())
	}
	defer resp.Body.Close()

	return nil
}
