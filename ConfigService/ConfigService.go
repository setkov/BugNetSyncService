package ConfigService

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type ConfigService struct {
	BugNetConnectionString string
	TfsBaseUri             string
	TfsАuthorizationToken  string
	IdleMode               bool
}

// New config service
func NewConfigService() *ConfigService {
	var config ConfigService
	if err := config.loadJson(); err != nil {
		log.Print("Error loadin configuration from file: ", err.Error())
	}
	config.loadEnvironment()
	return &config
}

// Load configuration from json file
func (c *ConfigService) loadJson() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&c); err != nil {
		return err
	}
	return nil
}

// Load configuration from environment variables
func (c *ConfigService) loadEnvironment() {
	if bugNetConnectionString, exists := os.LookupEnv("BUG_NET_CONNECTION_STRING"); exists {
		c.BugNetConnectionString = bugNetConnectionString
	}
	if tfsBaseUri, exists := os.LookupEnv("TFS_BASE_URI"); exists {
		c.TfsBaseUri = tfsBaseUri
	}
	if tfsАuthorizationToken, exists := os.LookupEnv("TFS_АUTHORIZATION_TOKEN"); exists {
		c.TfsАuthorizationToken = tfsАuthorizationToken
	}
	if idleMode, exists := os.LookupEnv("IDLE_MODE"); exists {
		c.IdleMode, _ = strconv.ParseBool(idleMode)
	}
}
