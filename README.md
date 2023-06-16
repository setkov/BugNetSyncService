## BugNetSyncService
Is a microservice for sync messages and attachments from BugNet HelpDesk system to TFS. Send status and error message to MSTeams or Telegram.

###### Sync sequence diagram:
![Sync sequence diagram](https://github.com/setkov/BugNetSyncService/blob/master/SequenceDiagram.png)

###### Configuraton:
set by a file in the root directory
config.json
```
{
  "BugNetConnectionString": "server=BugNet_Server;database=BugNet_Database;user id=BugNet_User;password=BugNet_User_Password;",
  "AttachmentServiceUrl": "support.bug-net.com/Attachment",
  "TfsBaseUri": "http://tfs_url/tfs/DefaultCollection/",
  "TfsАuthorizationToken": "tfs_token",
  "MSTeamsWebhookUrl": "MSteams webhook url",
  "TelegramToken": "token",
  "TelegramChatId": "chat_id",
  "IdleMode": false
}
```
or set environment variables (replace values from config file)
```
BUG_NET_CONNECTION_STRING - connection string to SQL server
BUG_NET_ATTACHMENT_SERVICE_URL - attachments servise url
TFS_BASE_URI - TFS base uri
TFS_АUTHORIZATION_TOKEN - TFS authorization token
MSTEAMS_WEBHOOK_URL - MSteams webhook url
TELEGRAM_TOKEN - Telegram token
TELEGRAM_CHAT_ID - Telegram chat id
IDLE_MODE - run service in idle mode (default value is false)
```    
###### Run service in docker
- get last docker image from [DockerHub](https://hub.docker.com/r/setkov/bug-net-sync-service/tags?page=1&ordering=last_updated)
- run the image using environment variables

###### View Message Queue
open in web browser link http://localhost:8080/