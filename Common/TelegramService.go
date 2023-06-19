package Common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TelegramService struct {
	Token  string
	ChatId string
	Sender string
}

type TelegramStatus struct {
	Ok          bool
	Error_code  int16
	Description string
}

// New TelegramService
func NewTelegramService(token string, chatId string, sender string) TelegramService {
	return TelegramService{
		Token:  token,
		ChatId: chatId,
		Sender: sender,
	}
}

// Send message
func (s TelegramService) SendMessage(message string) error {
	message = "<b>" + s.Sender + ":</b> " + message
	body, err := json.Marshal(map[string]interface{}{"chat_id": s.ChatId, "text": message, "parse_mode": "HTML"})
	if err != nil {
		return NewError("Prepare telegram message. " + err.Error())
	}

	uri := "https://api.telegram.org/bot" + s.Token + "/sendMessage"
	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return NewError("Post telegram message. " + err.Error())
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return NewError("Read request. " + err.Error())
	}

	var status TelegramStatus
	err = json.Unmarshal(bytes, &status)
	if err != nil {
		return NewError("Parse telegram status. " + err.Error())
	}

	if !status.Ok {
		return NewError(fmt.Sprintf("Send telegram message. %v %v", status.Error_code, status.Description))
	}

	return nil
}
