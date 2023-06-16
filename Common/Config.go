package Common

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

const configFileName string = "config.json"

type Config struct {
	ConnectionString      string
	AttachmentServiceUrl  string
	TfsBaseUri            string
	TfsАuthorizationToken string
	MSTeamsWebhookUrl     string
	TelegramToken         string
	TelegramChatId        string
	IdleMode              bool
}

// New config
func NewConfig() (*Config, error) {
	var config Config
	err := config.loadJsonFile()
	config.loadEnvironment()
	return &config, err
}

// Load configuration from json file
func (c *Config) loadJsonFile() error {
	if _, err := os.Stat(configFileName); errors.Is(err, os.ErrNotExist) {
		return NewWarning("Configuration file not found.")
	} else {
		file, err := os.Open(configFileName)
		if err != nil {
			return NewError("Open config file. " + err.Error())
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&c); err != nil {
			return NewError("Parse config file. " + err.Error())
		}
	}
	return nil
}

// Load configuration from environment variables
func (c *Config) loadEnvironment() {
	if bugNetConnectionString, exists := os.LookupEnv("BUG_NET_CONNECTION_STRING"); exists {
		c.ConnectionString = bugNetConnectionString
	}
	if attachmentServiceUrl, exists := os.LookupEnv("BUG_NET_ATTACHMENT_SERVICE_URL"); exists {
		c.AttachmentServiceUrl = attachmentServiceUrl
	}
	if tfsBaseUri, exists := os.LookupEnv("TFS_BASE_URI"); exists {
		c.TfsBaseUri = tfsBaseUri
	}
	if tfsАuthorizationToken, exists := os.LookupEnv("TFS_АUTHORIZATION_TOKEN"); exists {
		c.TfsАuthorizationToken = tfsАuthorizationToken
	}
	if msTeamsWebhookUrl, exists := os.LookupEnv("MSTEAMS_WEBHOOK_URL"); exists {
		c.MSTeamsWebhookUrl = msTeamsWebhookUrl
	}
	if telegramToken, exists := os.LookupEnv("TELEGRAM_TOKEN"); exists {
		c.TelegramToken = telegramToken
	}
	if telegramChatId, exists := os.LookupEnv("TELEGRAM_CHAT_ID"); exists {
		c.TelegramChatId = telegramChatId
	}
	if idleMode, exists := os.LookupEnv("IDLE_MODE"); exists {
		c.IdleMode, _ = strconv.ParseBool(idleMode)
	}
}
