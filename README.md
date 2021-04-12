## BugNetSyncService
BugNetSyncService is a microservice for sync message from BugNet HelpDesk system to TFS.

###### Sync sequence diagram:
![Sync sequence diagram](https://github.com/setkov/BugNetSyncService/blob/master/SequenceDiagram.png)

###### Configuraton:
configuration set by a file in the root directory
Config.json

  {

    "BugNetConnectionString": "server=BugNet_Server;database=BugNet_Database;user id=BugNet_User;password=BugNet_User_Password;",
    "TfsBaseUri": "http://tfs_url/tfs/DefaultCollection/",
    "TfsАuthorizationToken": "tfs_token",
    "IdleMode": false
  }
  
or from environment variables (replace values from config file)

    BUG_NET_CONNECTION_STRING - connection string to SQL server
    TFS_BASE_URI - TFS base uri
    TFS_АUTHORIZATION_TOKEN - TFS authorization token
    IDLE_MODE - run service in idle mode (default value is false)
    
###### Run service in docker
- get last docker image from [DockerHub](https://hub.docker.com/r/setkov/bug-net-sync-service/tags?page=1&ordering=last_updated)
- run the image using environment variables
