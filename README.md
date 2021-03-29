# BugNetSyncService

microservice for sync message from BugNet HelpDesk system to TFS

need configuration file in root directory 

Config.json

{
    
    "BugNetConnectionString": "server=BugNet_Server;database=BugNet_Database;user id=BugNet_User;password=BugNet_User_Password;",
    "TfsBaseUri": "http://tfs_url_/tfs/DefaultCollection/",
    "TfsАuthorizationToken": "tfs_token"

}
